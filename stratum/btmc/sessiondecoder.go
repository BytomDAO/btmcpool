package btmc

import (
	"errors"

	"github.com/segmentio/encoding/json"

	ss "github.com/bytom/btmcpool/stratum"
)

const (
	MLogin  string = "login"
	MSubmit string = "submit"
)

type btmDecoder struct{}

func NewBtmDecoder() ss.Decoder { return &btmDecoder{} }

func (d *btmDecoder) Decode(data []byte, session *ss.TcpSession) (ss.Request, error) {
	var req ss.JSONRpcReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	switch req.Method {
	case MLogin:
		return d.decodeLogin(data, session, req.Id)
	case MSubmit:
		return d.decodeSubmit(data, session, req.Id)
	default:
		return nil, errors.New("invalid data length")
	}
}

func (d *btmDecoder) decodeSubmit(data []byte, session *ss.TcpSession, id *json.RawMessage) (ss.Request, error) {
	var request submitRequest
	if err := json.Unmarshal(data, &request); err != nil {
		session.Error(request.Id, ss.ErrorUnknown)
		return nil, err
	}
	return NewSubmitRequest(request), nil
}

func (d *btmDecoder) decodeLogin(data []byte, session *ss.TcpSession, id *json.RawMessage) (ss.Request, error) {
	var request loginRequest
	if err := json.Unmarshal(data, &request); err != nil {
		session.Error(id, ss.ErrorUnknown)
		return nil, err
	}
	return NewLoginRequest(request), nil
}
