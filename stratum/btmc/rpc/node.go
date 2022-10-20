package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/encoding/json"

	"github.com/bytom/api"
	btmcrpc "github.com/bytom/blockchain/rpc"
	ss "github.com/bytom/btmcpool/stratum"
)

type BtmcClient struct {
	service string
	btmcCli *btmcrpc.Client
}

func NewBtmcClient(service string, nodeURL string) *BtmcClient {
	return &BtmcClient{
		service: service,
		btmcCli: &btmcrpc.Client{BaseURL: nodeURL},
	}
}

func (c *BtmcClient) GetWork() (*api.GetWorkResp, error) {
	var result ss.NodeJsonRpcResp
	if err := ss.CallWithMethod(c.service, "get-work", []string{}, &result); err != nil {
		return nil, err
	}

	if result.Data == nil {
		return nil, errors.New("empty reply for get work")
	}
	var reply api.GetWorkResp
	if err := json.Unmarshal(*result.Data, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func (c *BtmcClient) SubmitBlock(req interface{}) (interface{}, error) {
	var response = &api.Response{}
	c.btmcCli.Call(context.Background(), "/submit-work", req, response)

	switch response.Status {
	case api.FAIL:
		return nil, errors.New("error.bytomd")
	case "":
		return nil, errors.New("error.connect_bytomd")
	}
	return response.Data, nil
}

type Peer struct {
	RemoteAddr string `json:"remote_addr"`
	Height     int64  `json:"height"`
	Ping       string `json:"ping"`
}

func (c *BtmcClient) GetPeers() ([]*Peer, error) {
	var result ss.NodeJsonRpcResp
	if err := ss.CallWithMethod(c.service, "list-peers", []string{}, &result); err != nil {
		return nil, err
	}

	var peers []*Peer
	if err := json.Unmarshal(*result.Data, &peers); err != nil {
		return nil, err
	}

	return peers, nil
}

type balance struct {
	Amount int64 `json:"amount"`
}

func (c *BtmcClient) GetBalance() (int64, error) {
	var result ss.NodeJsonRpcResp
	if err := ss.CallWithMethod(c.service, "list-balances", []string{}, &result); err != nil {
		return 0, err
	}

	var balances []*balance
	if err := json.Unmarshal(*result.Data, &balances); err != nil {
		return 0, err
	}

	if len(balances) != 1 {
		return 0, fmt.Errorf("unexpected balance response")
	}

	return balances[0].Amount, nil
}
