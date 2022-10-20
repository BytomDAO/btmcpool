package superstratum

import (
	"sync"
)

// session Id manager monitors recycled session id.
// it assumes only one id getter (listener) and multiple recyclers (tcpSession.close()).
// use 2 bitsets to minimize session id creation overhead
type sessionIdManager struct {
	freeIds        []uint
	nextId         uint
	maxNumSessions uint
	mutex          sync.Mutex
}

func newSessionIdManager(maxNumSessions uint) *sessionIdManager {
	return &sessionIdManager{
		freeIds:        []uint{},
		maxNumSessions: maxNumSessions,
	}
}

func (s *sessionIdManager) getId() uint {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.freeIds) > 0 {
		id := s.freeIds[0]
		s.freeIds = s.freeIds[1:len(s.freeIds)]
		return id
	}

	id := s.nextId
	s.nextId++
	s.nextId = s.nextId % s.maxNumSessions
	return id
}

func (s *sessionIdManager) recycle(id uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.freeIds = append(s.freeIds, id)
}
