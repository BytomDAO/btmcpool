package datastruct

import (
	"sync/atomic"
)

// Thread-safe ring buffer implementation
// NOTE: items in the buffer are not guaranteed to be available in
// insertion order (i.e. an item in an earlier ring position may
// become available after an later item). Eventual consistency should follow.
type Ring struct {
	n int32
	p int32
	s []atomic.Value
}

func NewRing(n int32) *Ring {
	return &Ring{
		n: n,
		p: -1,
		s: make([]atomic.Value, n),
	}
}

func (r *Ring) Add(e interface{}) {
	p := atomic.AddInt32(&r.p, 1)
	r.s[p%r.n].Store(e)
}

func (r *Ring) Do(f func(interface{}) bool) {
	for _, av := range r.s {
		e := av.Load()
		if e != nil && f(e) {
			break
		}
	}
}
