package datastruct

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	r := NewRing(10)
	r.Add(1)
	r.Add(2)
	r.Add(3)
	sum := []int{0}
	r.Do(func(e interface{}) bool {
		sum[0] += e.(int)
		return false
	})
	assert.Equal(t, 6, sum[0])

	r.Add(4)
	r.Add(5)
	r.Add(6)
	sum[0] = 0
	r.Do(func(e interface{}) bool {
		sum[0] += e.(int)
		return false
	})
	assert.Equal(t, 21, sum[0])

	r.Add(7)
	r.Add(8)
	r.Add(9)
	r.Add(10)
	r.Add(11)
	r.Add(12)
	sum[0] = 0
	r.Do(func(e interface{}) bool {
		sum[0] += e.(int)
		return false
	})
	assert.Equal(t, 75, sum[0])
}

func TestStressRun(t *testing.T) {
	t.Skip("Skipping stress test in CI")

	add := func(r *Ring, d time.Duration, wg *sync.WaitGroup) {
		t.Log("add thread started")
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				wg.Done()
				return
			default:
				r.Add(rand.Int())
			}
		}
	}

	sum := func(r *Ring, d time.Duration, wg *sync.WaitGroup) {
		t.Log("sum thread started")
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				wg.Done()
				return
			default:
				sum := []int{0}
				r.Do(func(e interface{}) bool {
					sum[0] += e.(int)
					return false
				})
			}
		}
	}

	d := 60 * time.Second
	n := 4

	r := NewRing(10)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go add(r, d, &wg)
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go sum(r, d, &wg)
	}
	wg.Wait()
}
