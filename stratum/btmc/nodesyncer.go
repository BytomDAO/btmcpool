package btmc

import (
	"errors"
	"sync"

	"github.com/segmentio/encoding/json"

	"github.com/bytom/api"
	"github.com/bytom/btmcpool/common/logger"
	ss "github.com/bytom/btmcpool/stratum"
	"github.com/bytom/btmcpool/stratum/btmc/rpc"
)

type btmcNodeSyncer struct {
	client *rpc.BtmcClient
	bt     *api.GetWorkResp
	btLock sync.RWMutex

	latestHeight uint64
}

func NewBtmcNodeSyncer(service string, nodeURL string) (*btmcNodeSyncer, error) {
	return &btmcNodeSyncer{
		client:       rpc.NewBtmcClient(service, nodeURL),
		latestHeight: 0,
	}, nil
}

func (n *btmcNodeSyncer) fetchBlockTemplate() (ss.BlockTemplate, error) {
	reply, err := n.client.GetWork()
	if err != nil {
		return nil, err
	}

	header := reply.BlockHeader
	if header == nil {
		return nil, ErrNullBlockHeader
	}

	return &btmcBlockTemplate{
		version:                header.Version,
		height:                 header.Height,
		previousBlockHash:      &header.PreviousBlockHash,
		timestamp:              header.Time(),
		transactionsMerkleRoot: &header.TransactionsMerkleRoot,
		transactionStatusHash:  &header.TransactionStatusHash,
		nonce:                  header.Nonce,
		bits:                   header.Bits,
		seed:                   reply.Seed,
	}, nil
}

func (n *btmcNodeSyncer) Pull() (ss.BlockTemplate, error) {
	return n.fetchBlockTemplate()
}

func (n *btmcNodeSyncer) Submit(share ss.Share) error {
	btmcShare := share.(*btmcShare)
	rawdata, err := n.client.SubmitBlock(&api.SubmitWorkReq{BlockHeader: btmcShare.header})
	if err != nil {
		return err
	}

	resultrawdata, err := json.Marshal(rawdata)
	if err != nil {
		return err
	}
	var result bool
	if err := json.Unmarshal(resultrawdata, &result); err != nil {
		return err
	}
	if !result {
		logger.Error("block rejected", "nonce", btmcShare.nonce, "hash", btmcShare.blockHash)
		return nil
	}
	logger.Info("send nonce success", "nonce", btmcShare.nonce)
	return nil
}

func (n *btmcNodeSyncer) GetBt() (*api.GetWorkResp, error) {
	n.btLock.RLock()
	defer n.btLock.RUnlock()
	if n.bt == nil {
		return nil, errors.New("getting blocktemplate")
	}
	return n.bt, nil
}
