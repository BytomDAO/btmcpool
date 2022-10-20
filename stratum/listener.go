package superstratum

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/bytom/btmcpool/common/logger"
)

type listener struct {
	ctx            context.Context
	serverState    *ServerState
	port           int
	maxConn        int
	tcpListener    *net.TCPListener
	builder        *tcpSessionBuilder // builder for session
	diffBuilder    *diffAdjust        // builder for diffAdjust
	stopSignal     chan bool
	completeSignal chan bool
}

func newListener(
	ctx context.Context,
	port, maxConn int,
	serverState *ServerState,
	sessionBuilder *tcpSessionBuilder,
	diffBuilder *diffAdjust) (*listener, error) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	return &listener{
		ctx:            ctx,
		port:           port,
		maxConn:        maxConn,
		serverState:    serverState,
		tcpListener:    tcpListener,
		builder:        sessionBuilder,
		diffBuilder:    diffBuilder,
		stopSignal:     make(chan bool, 1),
		completeSignal: make(chan bool, 1),
	}, nil
}

// accepts new connection and creates session to handle
func (l *listener) listen() {
	logger.Info("start listening", "port", l.port)

	inflight := make(chan bool, l.maxConn)
	done := false

	for !done {
		conn, err := l.tcpListener.AcceptTCP()
		if err != nil {
			select {
			case <-l.stopSignal:
				done = true
			default:
				logger.Info("error accepting connection", "err", err)
			}
			continue
		}

		// close connection immediately if max # of connections is reached
		select {
		case inflight <- true:
		default:
			conn.Close()
			logger.Warn("maximum connections reached")
			continue
		}

		ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			conn.Close()
			logger.Info("error parsing remote ip", "err", err)
			continue
		}

		conn.SetKeepAlive(true)
		s, err := l.builder.build(l.ctx, conn, ip, l.diffBuilder)
		if err != nil {
			conn.Close()
			<-inflight
			logger.Error("error creating session", "session_ip", s.GetIp(), "err", err)
			continue
		}
		l.serverState.registerSession(s)

		logger.Info("connection accepted", "session_ip", s.GetIp(), "session_id", s.GetId())

		s.Run(func() error {
			s.dispatch(s.ctx)
			l.serverState.removeSession(s)
			<-inflight
			return errors.New("session closed")
		})
	}

	// clean up
	l.serverState.clearSessions()
	logger.Info("stop listening", "port", l.port)
	l.completeSignal <- true
}

// cleans up sessions and closes listener
func (l *listener) close() {
	l.stopSignal <- true
	l.tcpListener.Close()
	<-l.completeSignal
}
