package btmc

import (
	ss "github.com/bytom/btmcpool/stratum"
)

type LoginReq struct{ loginRequest }
type GetWorkReq struct{ getWorkRequest }
type SubmitReq struct{ submitRequest }
type KeepAlivedReq struct{ ss.JSONRpcReq }

func NewLoginRequest(req loginRequest) *LoginReq                   { return &LoginReq{req} }
func (r *LoginReq) Name() string                                   { return "login" }
func (r *LoginReq) Handle(session *ss.TcpSession) error            { return r.handleLogin(session) }
func (r *LoginReq) Forward(session *ss.TcpSession) error           { return nil }
func (r *LoginReq) CheckMiner(session *ss.TcpSession) (bool, bool) { return false, false }

func NewSubmitRequest(req submitRequest) *SubmitReq                 { return &SubmitReq{req} }
func (r *SubmitReq) Name() string                                   { return "submit" }
func (r *SubmitReq) Handle(session *ss.TcpSession) error            { return r.handleSubmit(session) }
func (r *SubmitReq) Forward(session *ss.TcpSession) error           { return nil }
func (r *SubmitReq) CheckMiner(session *ss.TcpSession) (bool, bool) { return false, false }
