package btmc

import (
	"math/big"
	"time"

	"github.com/bytom/protocol/bc"

	"github.com/bytom/btmcpool/common/mining/utils"
	ss "github.com/bytom/btmcpool/stratum"
	"github.com/bytom/btmcpool/stratum/btmc/util"
)

type btmcJob struct {
	id                     ss.JobId
	version                uint64
	height                 uint64
	previousBlockHash      *bc.Hash
	timestamp              time.Time
	transactionsMerkleRoot *bc.Hash
	transactionStatusHash  *bc.Hash
	bits                   uint64
	seed                   *bc.Hash
	nonce                  uint64
	diff                   *big.Int
}

func (j *btmcJob) GetId() ss.JobId {
	return j.id
}

func (j *btmcJob) GetDiff() uint64 {
	return j.diff.Uint64()
}

func (j *btmcJob) GetTarget() (string, bool, bool) {
	return "", false, false
}

func (j *btmcJob) Encode() (interface{}, error) {
	return ss.StratumJSONRpcNotify{
		Version: "2.0",
		Method:  "job",
		Params:  j.genReplyData(),
	}, nil
}

func (j *btmcJob) encodeLogin(login string) *jobReply {
	return &jobReply{
		Id:     login,
		Job:    j.genReplyData(),
		Status: "OK",
	}
}

func (j *btmcJob) genReplyData() *jobReplyData {
	return &jobReplyData{
		JobId:                  j.GetId().String(),
		Version:                utils.ToLittleEndianHex(j.version),
		Height:                 utils.ToLittleEndianHex(j.height),
		PreviousBlockHash:      j.previousBlockHash.String(),
		Timestamp:              utils.ToLittleEndianHex(uint64(j.timestamp.Unix())),
		TransactionsMerkleRoot: j.transactionsMerkleRoot.String(),
		TransactionStatusHash:  j.transactionStatusHash.String(),
		Nonce:                  utils.ToLittleEndianHex(uint64(j.nonce)),
		Bits:                   utils.ToLittleEndianHex(j.bits),
		Seed:                   j.seed.String(),
		Target:                 util.GetTargetHex(j.diff),
	}
}
