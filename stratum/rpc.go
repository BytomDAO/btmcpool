package superstratum

import (
	"bytes"
	"errors"
	"io"

	"github.com/segmentio/encoding/json"

	"github.com/bytom/btmcpool/common/rpc/http"
)

// bytom has a different json rpc format
// in order to use general 'Call' method, data & status were added to this struct
type NodeJsonRpcResp struct {
	Id     *json.RawMessage       `json:"id,omitempty"`
	Result *json.RawMessage       `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status string                 `json:"status,omitempty"`
	Data   *json.RawMessage       `json:"data,omitempty"`
}

func Call(service, method string, request interface{}, rpcResp *NodeJsonRpcResp) error {
	jsonReq, header := getHdAndReq(method, request)

	if err := http.CallImpl(service, "", header, jsonReq, rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return errors.New(rpcResp.Error["message"].(string))
	}
	return nil
}

func CallUrl(url, method string, request interface{}, rpcResp *NodeJsonRpcResp) error {
	jsonReq, header := getHdAndReq(method, request)

	var body io.Reader
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(jsonReq); err != nil {
		return err
	}
	body = &buf

	if err := http.SendRequest("POST", url, body, header, rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return errors.New(rpcResp.Error["message"].(string))
	}
	return nil
}

func CallWithMethod(service, method string, request interface{}, rpcResp *NodeJsonRpcResp) error {
	jsonReq, header := getHdAndReq(method, request)

	if err := http.CallImpl(service, method, header, jsonReq, rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return errors.New(rpcResp.Error["message"].(string))
	}
	return nil
}

func CallRawRequest(service, method string, request interface{}, rpcResp *NodeJsonRpcResp) error {
	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	if err := http.CallImpl(service, method, header, request, rpcResp); err != nil {
		return err
	}
	if rpcResp.Error != nil {
		return errors.New(rpcResp.Error["message"].(string))
	}
	return nil
}

func getHdAndReq(method string, request interface{}) (map[string]interface{}, map[string]string) {
	jsonReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  request,
		"id":      0,
	}
	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	return jsonReq, header
}
