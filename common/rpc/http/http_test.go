package http

import (
	"fmt"
	h "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bytom/btmcpool/common/rpc/hostprovider"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	ts := httptest.NewServer(h.HandlerFunc(func(w h.ResponseWriter, r *h.Request) {
		fmt.Fprintln(w, "{\"Method\": \""+r.URL.Path+"\", \"Value\":16}")
	}))
	defer ts.Close()

	hostprovider.InitStaticProvider(
		map[string][]string{
			"test_service": []string{ts.URL},
		})

	Init(time.Second)

	type Result struct {
		Method string
		Value  int
	}

	var result Result
	assert.NoError(t, Call("test_service", "method1", "", &result))
	assert.Equal(t, "/method1", result.Method)
	assert.Equal(t, 16, result.Value)

	assert.NoError(t, Call("test_service", "method2", "", &result))
	assert.Equal(t, "/method2", result.Method)
	assert.Equal(t, 16, result.Value)

	assert.Error(t, Call("test_service_nonexist", "method", "", &result))
}
