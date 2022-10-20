package btmc

import (
	"github.com/segmentio/encoding/json"

	ss "github.com/bytom/btmcpool/stratum"
)

type loginRequest struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params struct {
		Login string `json:"login"`
		Pass  string `json:"pass"`
		Agent string `json:"agent"`
	} `json:"params"`
}

type loginResp struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  jobReply         `json:"result"`
	Error   *ss.ErrorReply   `json:"error"`
}

type getWorkRequest struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params struct {
		Id string `json:"id"`
	} `json:"params"`
}

type submitRequest struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params struct {
		Id     string `json:"id"`
		JobId  string `json:"job_id"`
		Nonce  string `json:"nonce"`
		Result string `json:"result"`
	} `json:"params"`
}

type submitResp struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  statusReply      `json:"result"`
	Error   *ss.ErrorReply   `json:"error"`
}

type f2poolResp struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  bool             `json:"result"`
	Error   *ss.ErrorReply   `json:"error"`
}

type jobReply struct {
	Id     string        `json:"id"`
	Job    *jobReplyData `json:"job"`
	Status string        `json:"status"`
}

type jobReplyData struct {
	Version                string `json:"version"`
	Height                 string `json:"height"`
	PreviousBlockHash      string `json:"previous_block_hash"`
	Timestamp              string `json:"timestamp"`
	TransactionsMerkleRoot string `json:"transactions_merkle_root"`
	TransactionStatusHash  string `json:"transaction_status_hash"`
	Nonce                  string `json:"nonce"`
	Bits                   string `json:"bits"`
	JobId                  string `json:"job_id"`
	Seed                   string `json:"seed"`
	Target                 string `json:"target"`
}

type statusReply struct {
	Status string `json:"status"`
}
