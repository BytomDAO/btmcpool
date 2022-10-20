package btmc

import (
	"errors"
	"math/big"

	"github.com/bytom/protocol/bc"
	"github.com/bytom/protocol/bc/types"

	"github.com/bytom/btmcpool/common/logger"
	ss "github.com/bytom/btmcpool/stratum"
)

const (
	verMagicNum   = uint64(1)
	defaultReward = uint64(41250000000)
	defaultFee    = uint64(300)
)

type btmcShare struct {
	job    *btmcJob
	worker *ss.Worker

	nonce     uint64
	result    string
	header    *types.BlockHeader
	blockHash *bc.Hash
	netDiff   *big.Int

	state  ss.ShareState
	reason ss.RejectReason
}

// build block from the share for node submission
func (s *btmcShare) BuildBlock() (ss.BlockTemplate, error) {
	// not implemented
	logger.Fatal("BuildBlock not implemented")
	return nil, nil
}

// build pb sharelog from the share for logging
func (s *btmcShare) BuildLog(port uint64) ([]byte, error) {
	return nil, errors.New("not support")
}

// update share state
func (s *btmcShare) UpdateState(state ss.ShareState, reason ss.RejectReason) error {
	s.state = state
	s.reason = reason
	return nil
}

func (s *btmcShare) GetState() ss.ShareState {
	return s.state
}

func (s *btmcShare) GetReason() ss.RejectReason {
	return s.reason
}

func (s *btmcShare) GetWorker() *ss.Worker {
	return s.worker
}
