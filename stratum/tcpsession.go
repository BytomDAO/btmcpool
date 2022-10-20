package superstratum

import (
	"bufio"
	"context"
	"io"
	"math/big"
	"net"
	"time"

	"github.com/segmentio/encoding/json"
	"golang.org/x/sync/errgroup"

	"github.com/bytom/btmcpool/common/datastruct"
	"github.com/bytom/btmcpool/common/logger"
)

const (
	maxReqSize = 10 * 1024
)

// SessionState codes
type SessionState int

// valid session states
const (
	SStateUndefined  SessionState = iota // undefined session state (default)
	SStateConnected                      // session is connected
	SStateSubscribed                     // session is subscribed
	SStateAuthorized                     // session is authorized
	SStateClosed                         // session is closed
)

// TcpSessionBuilder is the builder for tcp session
type tcpSessionBuilder struct {
	state    *ServerState
	syncer   NodeSyncer
	verifier Verifier
	timeout  time.Duration
	interval time.Duration
	decoder  Decoder

	dataBuilder SessionDataBuilder
}

//NewTcpSessionBuilder creates a TcpSessionBuilder
func newTcpSessionBuilder(state *ServerState,
	syncer NodeSyncer,
	verifier Verifier,
	timeout time.Duration,
	interval time.Duration,
	dataBuilder SessionDataBuilder,
	decoder Decoder) *tcpSessionBuilder {
	return &tcpSessionBuilder{
		state:       state,
		syncer:      syncer,
		verifier:    verifier,
		timeout:     timeout,
		interval:    interval,
		decoder:     decoder,
		dataBuilder: dataBuilder,
	}
}

// Build builds a TcpSession
func (b *tcpSessionBuilder) build(ctx context.Context, conn *net.TCPConn, ip string, diffBuilder *diffAdjust) (*TcpSession, error) {
	id := b.state.idManager.getId()
	return newTcpSession(ctx, id, conn, ip, b.state, b.syncer, b.verifier, b.timeout, b.interval, b.dataBuilder.Build(id), diffBuilder, b.decoder), nil
}

// TcpSession is a session implementation with underlying TCP connection
// This is intended to be used as base (coin-independent) session
// Individual coin session should extend TcpSession with its own internal data
type TcpSession struct {
	id           uint // unique session id
	ip           string
	conn         *net.TCPConn
	Encoder      *json.Encoder
	timeout      time.Duration // connection timeout
	serverState  *ServerState  // global server state
	notifyMehtod string

	diffAdjust *diffAdjust
	SessionCtl *SessionCtl

	// services needed by handlers
	syncer  NodeSyncer
	verifer Verifier
	state   SessionState

	jobInterval time.Duration    // interval of sending new jobs to workers. zero means no disable scheduling
	jobHist     *datastruct.Ring // outstanding jobs, job id as key
	jobSignal   chan bool

	sessionData    SessionData
	group          *errgroup.Group
	ctx            context.Context
	CancelSchedule context.CancelFunc // should be changed to map with the corresponding ctx
	decoder        Decoder
}

func newTcpSession(
	ctx context.Context,
	id uint,
	conn *net.TCPConn,
	ip string,
	serverState *ServerState,
	syncer NodeSyncer,
	verifier Verifier,
	timeout time.Duration,
	interval time.Duration,
	data SessionData,
	diffBuilder *diffAdjust,
	decoder Decoder,
) *TcpSession {
	session := &TcpSession{
		id:          id,
		ip:          ip,
		conn:        conn,
		Encoder:     json.NewEncoder(conn),
		serverState: serverState,
		timeout:     timeout,
		syncer:      syncer,
		verifer:     verifier,
		jobInterval: interval,
		state:       SStateConnected,
		decoder:     decoder,
		jobHist:     datastruct.NewRing(8),
		jobSignal:   make(chan bool, 1),
		sessionData: data,
	}

	session.group, session.ctx = errgroup.WithContext(ctx)

	session.diffAdjust = diffBuilder
	session.SessionCtl = NewSessionCtl()

	var scheduleCtx context.Context
	scheduleCtx, session.CancelSchedule = context.WithCancel(session.ctx)
	session.Run(func() error { return session.scheduleJobs(scheduleCtx) })
	return session
}

func (s *TcpSession) GetState() SessionState {
	return s.state
}

func (s *TcpSession) SetState(state SessionState) {
	s.state = state
}

func (s *TcpSession) GetSessionData() SessionData {
	return s.sessionData
}

// record a new job
func (s *TcpSession) AddJob(job Job) error {
	s.jobHist.Add(job)

	return nil
}

func (s *TcpSession) FindJob(id JobId) Job {
	var job Job
	s.jobHist.Do(func(e interface{}) bool {
		if e.(Job).GetId() == id {
			job = e.(Job)
			return true
		}
		return false
	})
	return job
}

// clear all jobs
func (s *TcpSession) ClearJobs() {
}

func (s *TcpSession) GetJob() (Job, error) {
	blockTemplate := s.serverState.GetBlockTemplate()
	if blockTemplate == nil {
		return nil, ErrGetBlockTemplate
	}
	return blockTemplate.CreateJob(s)
}

func (s *TcpSession) GetId() uint {
	return s.id
}

func (s *TcpSession) GetIp() string {
	return s.ip
}

func (s *TcpSession) SendJob() {
	// alway non-blocking
	select {
	case s.jobSignal <- true:
	default:
	}
}

func (s *TcpSession) GetServerState() *ServerState {
	return s.serverState
}

func (s *TcpSession) dispatch(ctx context.Context) {
	logger.Info("session dispatch", "session_id", s.GetId(), "session_ip", s.GetIp())
	connbuff := bufio.NewReaderSize(s.conn, maxReqSize)
	defer s.close()

	var (
		data     []byte
		isPrefix bool
		err      error
	)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, isPrefix, err = connbuff.ReadLine()
			if isPrefix {
				logger.Error("socket flood detected", "session_id", s.id, "session_ip", s.ip)
				return
			} else if err == io.EOF {
				logger.Error("client disconnected", "session_id", s.id, "session_ip", s.ip)
				return
			} else if err != nil {
				// filter part ot ip addresses so as to avoid reporting connection reset
				_, ipv4Net, ipErr := net.ParseCIDR(s.ip + "/10")
				if ipErr != nil {
					logger.Error("parse cidr err", "session_id", s.id, "session_ip", s.ip, "err", err)
				} else if ipv4Net.String() != "100.64.0.0/10" {
					logger.Error("error reading", "session_id", s.id, "session_ip", s.ip, "err", err)
				}
				return
			}

			if s.decoder != nil {
				req, err := s.decoder.Decode(data, s)
				if err != nil {
					logger.Error("fail to decode ", "session_id", s.id, "session_ip", s.ip, "err", err, "data", data)
					return
				}
				// reset idle timeout
				s.setDeadline()

				inBanlist, needClose := req.CheckMiner(s)
				if needClose {
					return
				}

				if inBanlist && s.serverState.GetConnCtrl().fwdTrigger {
					err = req.Forward(s)
				} else {
					err = req.Handle(s)
				}

				if err != nil {
					// bail out when handler returns error
					if err != ErrBannedMiner {
						logger.Error("handler error",
							"session_id", s.id,
							"session_ip", s.ip,
							"handler", req.Name(),
							"error", err)
					}
					return
				}
			} else {
				if len(data) > 1 {
					var req JSONRpcReq
					err = json.Unmarshal(data, &req)
					if err != nil {
						logger.Error("malformed request", "session_id", s.id, "session_ip", s.ip, "err", err, "data", string(data))
						return
					}
					// reset idle timeout
					s.setDeadline()

					handler, ok := s.serverState.handlerMap[req.Method]
					if !ok {
						logger.Error("unsupported method", "method", req.Method, "session_ip", s.GetIp(), "session_id", s.GetId())
						s.Error(req.Id, ErrorUnsupported)
						continue
					}

					if err = handler(s, data); err != nil {
						// bail out when handler returns error
						if err != ErrBannedMiner {
							logger.Error("handler error",
								"session_id", s.id,
								"session_ip", s.ip,
								"method", req.Method,
								"error", err)
						}
						return
					}
				}
			}
		}
	}
}

// Notify sends server notification to client
func (s *TcpSession) Notify(message interface{}) error {
	if err := s.Encoder.Encode(&message); err != nil {
		return err
	}
	// reset idle timeout
	s.setDeadline()
	return nil
}

// Reply sends reply to client request
func (s *TcpSession) Reply(id *json.RawMessage, result interface{}) error {
	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
	if err := s.Encoder.Encode(&message); err != nil {
		return err
	}

	return nil
}

// ReplyOmitError sends reply to client request omitting empty error
func (s *TcpSession) ReplyOmitError(id *json.RawMessage, result interface{}) error {
	message := JSONRpcOmitErrorResp{Id: id, Version: "2.0", Result: result}
	if err := s.Encoder.Encode(&message); err != nil {
		return err
	}
	return nil
}

// SetTarget sends target to client request
func (s *TcpSession) SetTarget(params interface{}) error {
	message := StratumJSONRpcNotify{Version: "2.0", Method: "mining.set_target", Params: params}
	if err := s.Encoder.Encode(&message); err != nil {
		return err
	}
	return nil
}

// SetDiff sends diff to client request
func (s *TcpSession) SetDiff(params interface{}) error {
	message := StratumJSONRpcNotify{Version: "2.0", Method: "mining.set_difficulty", Params: params}
	if err := s.Encoder.Encode(&message); err != nil {
		return err
	}
	return nil
}

// Error sends error message to client reqeust
func (s *TcpSession) Error(id *json.RawMessage, errType ErrorType) error {
	message := JSONRpcResp{Id: id, Version: "2.0", Error: errType.genError()}
	if err := s.Encoder.Encode(&message); err != nil {
		logger.Warn("send error",
			"session_id", s.GetId(),
			"session_ip", s.GetIp(),
			"err_type", errType.genError().Message,
			"send_err", err)
		return err
	}
	return nil
}

// GetDiff returns calculated difficulty
func (s *TcpSession) GetDiff() *big.Int {
	return s.diffAdjust.GetDiff()
}

func (s *TcpSession) GetVerifier() Verifier {
	return s.verifer
}

func (s *TcpSession) GetNodeSyncer() NodeSyncer {
	return s.syncer
}

// manually close net conn to recycle the go routine
func (s *TcpSession) CloseTcpConn() {
	s.conn.Close()
}

func (s *TcpSession) close() {
	s.conn.Close()
	s.CancelSchedule() // send stop signal to job scheduler
	s.state = SStateClosed
	s.serverState.idManager.recycle(s.id)
	logger.Info("session closed", "session_id", s.GetId(), "session_ip", s.GetIp())
}

// set connection timeout
func (s *TcpSession) setDeadline() {
	// connection timeout does not automatically reset,
	// meaning connection will timeout even if there are activities going on.
	//
	// currently, 2 places we do timeout reset:
	// 1. when session dispatcher receives a new message from client
	// 2. when session job scheduler notifies a new job to the client
	s.conn.SetDeadline(time.Now().Add(s.timeout))
}

// schedule repeated job notifications
func (s *TcpSession) scheduleJobs(ctx context.Context) error {
	logger.Info("start job scheduler", "session_id", s.id)
	var ticker *time.Ticker
	if s.jobInterval <= 0 {
		ticker = time.NewTicker(1 * time.Minute)
		// stop the ticker right away. The channel is still open, but never fire
		ticker.Stop()
	} else {
		ticker = time.NewTicker(s.jobInterval)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		case <-s.jobSignal:
		}
		if s.state != SStateAuthorized {
			// only send job notify when session is authorized
			continue
		}
		//logger.Info("notifiying new job", "session_id", s.GetId(), "session_ip", s.GetIp())
		job, err := s.GetJob()
		if err != nil {
			logger.Error("fail to create job",
				"session_id", s.GetId(),
				"session_ip", s.GetIp(),
				"error", err)
			continue
		}
		// rpc msg to miner: set_target or set_difficulty, needed by some coins
		// Now zec/zen/ckb use set_target,but close
		// meter use set_difficulty,open
		// other coins diff was set within job together
		target, setTarget, setDiff := job.GetTarget()
		if setTarget {
			if err := s.SetTarget([]interface{}{target}); err != nil {
				logger.Error("fail to notify target",
					"session_id", s.GetId(),
					"session_ip", s.GetIp(),
					"job_id", job.GetId(),
					"error", err)
				continue
			}
		}
		if setDiff {
			if err := s.SetDiff([]interface{}{job.GetDiff()}); err != nil {
				logger.Error("fail to notify difficulty",
					"session_id", s.GetId(),
					"session_ip", s.GetIp(),
					"job_id", job.GetId(),
					"error", err)
				continue
			}
		}
		message, err := job.Encode()
		if err != nil {
			logger.Error("fail to encode job",
				"session_id", s.GetId(),
				"session_ip", s.GetIp(),
				"job_id", job.GetId(),
				"error", err)
			continue
		}

		if err := s.Notify(message); err != nil {
			logger.Error("fail to notify",
				"session_id", s.GetId(),
				"session_ip", s.GetIp(),
				"job_id", job.GetId(),
				"error", err)
			continue
		}
		if err := s.AddJob(job); err != nil {
			logger.Error("fail to add job",
				"session_id", s.GetId(),
				"session_ip", s.GetIp(),
				"job_id", job.GetId(),
				"error", err)
		}
	}
}

func (s *TcpSession) Run(f func() error) {
	s.group.Go(func() error {
		return f()
	})
}
