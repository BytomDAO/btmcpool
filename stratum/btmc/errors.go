package btmc

import "errors"

var (
	// ErrNullBlockHeader means block header is null
	ErrNullBlockHeader = errors.New("block header is null")
	// ErrStaleShare means height is not the same as block height
	ErrStaleShare = errors.New("error.stale share")
	// ErrBannedMiner
	ErrBannedMiner = errors.New("error.banned miner")
	ErrCloseSignal = errors.New("context close signal")
	ErrNilWire     = errors.New("error.nil wire")
)
