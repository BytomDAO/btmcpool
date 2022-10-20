package static

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	// create new sdclient with mysql backend
	c := NewProvider(
		"test_client",
		[]string{
			"192.168.0.0:80",
			"192.168.0.1:81",
			"192.168.0.2:82",
		},
	)

	// query
	servers, err := c.Get()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(servers))

	m := map[string]bool{}
	for _, s := range servers {
		m[s] = true
	}
	assert.True(t, m["192.168.0.0:80"])
	assert.True(t, m["192.168.0.1:81"])
	assert.True(t, m["192.168.0.2:82"])
}
