package superstratum

import "errors"

var (
	// ErrSubscribeParams means subscribe params are illegal
	ErrSubscribeParams = errors.New("malformed subscribe params")
	// ErrLoginParams means login params are illegal
	ErrLoginParams = errors.New("malformed login params")
	// ErrConfigureParams means submit params are illegal
	ErrConfigureParams = errors.New("malformed configure params")
	// ErrSubmitParams means submit params are illegal
	ErrSubmitParams = errors.New("malformed submit params")
	// ErrNotRegistered means submitted worker has not registered
	ErrNotRegistered = errors.New("worker not registered")
	// ErrGetWorkParams means getwork params are illegal
	ErrGetWorkParams = errors.New("malformed get-work params")
	// ErrGetBlockTemplate means can not get blockTemplate from serverState
	ErrGetBlockTemplate = errors.New("can't get block template")
	ErrBannedMiner      = errors.New("banned miner access")
	ErrBannedIp         = errors.New("banned ip access")
	ErrEmptyWire        = errors.New("empty wire")
)
