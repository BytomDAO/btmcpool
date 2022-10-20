package btmc

import (
	"math/big"

	"github.com/bytom/consensus/difficulty"
	algorithm "github.com/bytom/mining/tensority/go_algorithm"
	"github.com/bytom/protocol/bc/types"

	ss "github.com/bytom/btmcpool/stratum"
)

type btmcVerifier struct {
	serverState *ss.ServerState
}

func NewBtmcVerifier(state *ss.ServerState) (*btmcVerifier, error) {
	return &btmcVerifier{
		serverState: state,
	}, nil
}

func (v *btmcVerifier) Verify(share ss.Share) error {
	btmcShare := share.(*btmcShare)
	btmcJob := btmcShare.job
	btmcShare.header = &types.BlockHeader{
		Version:           btmcJob.version,
		Height:            btmcJob.height,
		PreviousBlockHash: *btmcJob.previousBlockHash,
		Timestamp:         uint64(btmcJob.timestamp.Unix()),
		BlockCommitment: types.BlockCommitment{
			TransactionsMerkleRoot: *btmcJob.transactionsMerkleRoot,
			TransactionStatusHash:  *btmcJob.transactionStatusHash,
		},
		Nonce: btmcShare.nonce,
		Bits:  btmcJob.bits,
	}
	shareHeader := btmcShare.header
	headerHash := shareHeader.Hash()
	cmpHash := algorithm.LegacyAlgorithm(&headerHash, btmcJob.seed)
	if cmpHash == nil {
		share.UpdateState(ss.ShareStateRejected, ss.RejectReasonUndefined)
		return nil
	}

	btmcShare.blockHash = &headerHash
	bMax := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)
	bits := difficulty.BigToCompact(big.NewInt(0).Div(bMax, btmcShare.netDiff))
	if difficulty.HashToBig(cmpHash).Cmp(difficulty.CompactToBig(bits)) <= 0 {
		share.UpdateState(ss.ShareStateBlock, ss.RejectReasonPass)
		return nil
	}

	shareBits := difficulty.BigToCompact(big.NewInt(0).Div(bMax, btmcJob.diff))
	if difficulty.HashToBig(cmpHash).Cmp(difficulty.CompactToBig(shareBits)) > 0 {
		share.UpdateState(ss.ShareStateRejected, ss.RejectReasonLowDiff)
		return nil
	}
	share.UpdateState(ss.ShareStateAccepted, ss.RejectReasonPass)
	return nil
}
