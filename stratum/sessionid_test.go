package superstratum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionIdMgr(t *testing.T) {
	mgr := newSessionIdManager(10)
	id := mgr.getId()
	assert.Equal(t, uint(0), id)
	id = mgr.getId()
	assert.Equal(t, uint(1), id)
	id = mgr.getId()
	assert.Equal(t, uint(2), id)
	id = mgr.getId()
	assert.Equal(t, uint(3), id)
	id = mgr.getId()
	assert.Equal(t, uint(4), id)
	id = mgr.getId()
	assert.Equal(t, uint(5), id)

	mgr.recycle(2)
	mgr.recycle(3)
	id = mgr.getId()
	assert.Equal(t, uint(2), id)
	mgr.recycle(2)
	id = mgr.getId()
	assert.Equal(t, uint(3), id)
	id = mgr.getId()
	assert.Equal(t, uint(2), id)
	id = mgr.getId()
	assert.Equal(t, uint(6), id)
}
