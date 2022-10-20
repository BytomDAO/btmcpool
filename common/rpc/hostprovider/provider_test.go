package hostprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatic(t *testing.T) {
	err := InitStaticProvider(
		map[string][]string{
			"service1": []string{
				"192.168.1.1:81",
				"192.168.1.2:81",
				"192.168.1.3:81",
			},
			"service2": []string{
				"192.168.2.1:82",
				"192.168.2.2:82",
			},
		})
	assert.NoError(t, err)
	addrs, err := Get("service1")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(addrs))
	addrs, err = Get("service2")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(addrs))
}
