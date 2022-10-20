package superstratum

import (
	"context"
	"sync"
	"time"

	"github.com/cornelk/hashmap"

	"github.com/bytom/btmcpool/common/logger"
)

// HandlerFunc signature
// ONLY return error when it is neccessary to close the session
// params:
// current session pointer
// received bytes
type HandlerFunc func(*TcpSession, []byte) error

// ServerState stores server level shared state
type ServerState struct {
	ctx           context.Context
	serverId      int // used for distinguish from peer servers
	handlerMap    map[string]HandlerFunc
	templateLock  sync.RWMutex
	blockTemplate BlockTemplate
	sessionMap    *hashmap.HashMap // active sessions. Must allow concurrent access.
	idManager     *sessionIdManager
	ssConnCtrl    *connCtrl
	BtSignal      chan bool
}

// InitServerState initializes server states
func InitServerState(ctx context.Context, connCtrl *connCtrl, id int, maxSessions uint) (*ServerState, error) {
	logger.Info("init server state", "id", id)
	return &ServerState{
		ctx:        ctx,
		serverId:   id,
		handlerMap: map[string]HandlerFunc{},
		sessionMap: &hashmap.HashMap{},
		idManager:  newSessionIdManager(maxSessions),
		ssConnCtrl: connCtrl,
		BtSignal:   make(chan bool, 1),
	}, nil
}

func (s *ServerState) GetId() int {
	return s.serverId
}

// RegisterHandler registers a new JSON method handler
func (s *ServerState) RegisterHandler(method string, handler HandlerFunc) error {
	if _, ok := s.handlerMap[method]; ok {
		logger.Warn("overwriting existing handler", "method", method)
	}
	s.handlerMap[method] = handler
	return nil
}

func (s *ServerState) GetConnCtrl() *connCtrl {
	return s.ssConnCtrl
}

func (s *ServerState) registerSession(sess *TcpSession) {
	s.sessionMap.Set(sess.GetId(), sess)
}

func (s *ServerState) removeSession(sess *TcpSession) {
	s.sessionMap.Del(sess.GetId())
	// close the session if neccessary
	if sess.GetState() != SStateClosed {
		sess.close()
	}
}

func (s *ServerState) clearSessions() {
	for m := range s.sessionMap.Iter() {
		m.Value.(*TcpSession).close()
	}
}

func (s *ServerState) updateBlockTemplate(template BlockTemplate) bool {
	s.templateLock.Lock()
	defer s.templateLock.Unlock()
	if s.blockTemplate == nil || s.blockTemplate.Compare(template) < 0 {
		s.blockTemplate = template
		return true
	}
	return false
}

func (s *ServerState) GetBlockTemplate() BlockTemplate {
	s.templateLock.RLock()
	defer s.templateLock.RUnlock()
	return s.blockTemplate
}

func (s *ServerState) broadcast() {
	template := s.GetBlockTemplate()
	if template == nil {
		return
	}

	logger.Info("broadcasting new jobs", "sessions", s.sessionMap.Len())

	n, bcast := 0, make(chan int, 1024*16)
	for m := range s.sessionMap.Iter() {
		n++
		bcast <- n
		go func(sess *TcpSession) {
			// current jobs are expired now, clear them all
			sess.ClearJobs()
			sess.SendJob()
			<-bcast
		}(m.Value.(*TcpSession))
	}
}

// NewServer connects different server conponents together and starts listening to new connections
func NewServer(port, maxConn int,
	state *ServerState,
	syncer NodeSyncer,
	syncInterval time.Duration,
	verifier Verifier,
	timeout time.Duration,
	interval time.Duration,
	dataBuilder SessionDataBuilder,
	diffBuilder *diffAdjust,
	decoder Decoder) error {
	if syncer != nil {
		go startNodeSyncerTicker(state, syncer, syncInterval)
	}

	sessionBuilder := newTcpSessionBuilder(state, syncer, verifier, timeout, interval, dataBuilder, decoder)
	listener, err := newListener(state.ctx, port, maxConn, state, sessionBuilder, diffBuilder)
	if err != nil {
		return err
	}
	go listener.listen()

	go state.ssConnCtrl.scanning(1 * time.Hour)
	return nil
}
