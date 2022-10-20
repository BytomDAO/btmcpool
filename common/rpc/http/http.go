package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/bitly/go-hostpool"
	"github.com/segmentio/encoding/json"

	"github.com/bytom/btmcpool/common/rpc/hostprovider"
)

// constants for long HTTP connections
const (
	MaxIdleConnections int = 30
)

type client struct {
	sync.RWMutex
	pools      map[string]hostpool.HostPool
	httpClient *http.Client
}

var c *client

func Call(service, method string, request, result interface{}) error {
	return c.callImpl(service, method, nil, request, result)
}

// CallJson calls a remote procedure on another node like eth/sipc json format
func CallImpl(service, method string, header map[string]string, request, result interface{}) error {
	return c.callImpl(service, method, header, request, result)
}

func Init(timeout time.Duration) {
	c = &client{
		pools: map[string]hostpool.HostPool{},
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: MaxIdleConnections,
			},
		},
	}
}

func (c *client) getHost(service string) (hostpool.HostPoolResponse, error) {
	c.RLock()
	p := c.pools[service]
	c.RUnlock()

	// fast path
	if p != nil {
		return p.Get(), nil
	}

	// slow path
	hosts, err := hostprovider.Get(service)
	if err != nil {
		return nil, err
	}
	newPool := hostpool.NewEpsilonGreedy(hosts, 0, &hostpool.LinearEpsilonValueCalculator{})

	c.Lock()
	defer c.Unlock()

	p = c.pools[service]
	if p != nil {
		return p.Get(), nil
	}

	c.pools[service] = newPool

	return newPool.Get(), nil
}

// Call calls a remote procedure on another node
func (c *client) callImpl(service, method string, header map[string]string, request, result interface{}) error {
	var body io.Reader
	if request != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(request); err != nil {
			return err
		}
		body = &buf
	}

	h, err := c.getHost(service)
	if err != nil {
		return err
	}
	url := h.Host()
	if len(method) > 0 {
		url += "/" + method
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		// We don't mark host as failure for Request construction failure
		return err
	}
	if header != nil {
		// set header fields for json data format request
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
			h.Mark(err)
		}
		return err
	}
	h.Mark(nil)

	if resp != nil {
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func ReadUrl(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &result); err != nil {
		return fmt.Errorf("error decoding address info, resp %v, err %v", string(bytes), err)
	}

	return nil
}

func SendRequest(method, url string, reader io.Reader, headers map[string]string, result interface{}, timeout ...time.Duration) error {
	expireTime := time.Second
	if len(timeout) > 0 {
		expireTime = timeout[0]
	}
	c := &http.Client{
		Timeout: expireTime,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
	}
	return sendRequest(c, method, url, reader, headers, result)
}

func SendRequestWithCli(c *http.Client, method, url string, reader io.Reader, headers map[string]string, result interface{}) error {
	return sendRequest(c, method, url, reader, headers, result)
}

func sendRequest(cli *http.Client, method, url string, reader io.Reader, headers map[string]string, result interface{}) error {
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := cli.Do(request)
	if err != nil {
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
		}
		if result != nil {
			return json.NewDecoder(resp.Body).Decode(result)
		}
	}
	return nil
}
