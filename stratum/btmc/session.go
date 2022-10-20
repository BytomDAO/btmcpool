package btmc

import (
	"math"

	"github.com/cornelk/hashmap"

	ss "github.com/bytom/btmcpool/stratum"
)

const serverIdOffset = uint(60)

type btmcSessionData struct {
	nonce    uint64
	worker   *ss.Worker
	submitId *hashmap.HashMap
}

func (s *btmcSessionData) GetWorker() *ss.Worker {
	return s.worker
}

func (s *btmcSessionData) SetWorker(worker *ss.Worker) {
	s.worker = worker
}

func (s *btmcSessionData) getNonce() uint64 {
	return s.nonce
}

type btmcSessionDataBuilder struct {
	id              uint64
	maxSessions     int
	sessionIdOffset uint
}

func NewBtmcSessionDataBuilder(serverId uint64, maxSessions int) *btmcSessionDataBuilder {
	sessionIdOffset := serverIdOffset - uint(math.Ceil(math.Log2(float64(maxSessions))))
	if sessionIdOffset == serverIdOffset {
		sessionIdOffset--
	}
	return &btmcSessionDataBuilder{
		id:              serverId,
		maxSessions:     maxSessions,
		sessionIdOffset: sessionIdOffset,
	}
}

// Build builds a btmSession
func (b *btmcSessionDataBuilder) Build(sessionId uint) ss.SessionData {
	return &btmcSessionData{
		nonce:    (b.id << serverIdOffset) | (uint64(sessionId) << b.sessionIdOffset),
		worker:   nil,
		submitId: &hashmap.HashMap{},
	}
}
