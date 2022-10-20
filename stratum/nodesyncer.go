package superstratum

import (
	"time"

	"github.com/bytom/btmcpool/common/logger"
)

// NodeSyncer defines common interface for block syncronization with node
type NodeSyncer interface {
	// pull latest block info from node
	// return nil if no new block is available
	Pull() (BlockTemplate, error)

	// submit accepted share(block) to node
	Submit(Share) error
}

// startNodeSyncerTicker start node syncer and keep polling for new block
func startNodeSyncerTicker(state *ServerState, syncer NodeSyncer, interval time.Duration) {
	logger.Info("start node syncing", "interval", interval)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		if err := broadcastBt(state, syncer); err != nil {
			logger.Error("node syncing", "error", err)
		}
	}
}

// startNodeSyncerSignal start node syncer and keep polling for new block
func startNodeSyncerSignal(state *ServerState, syncer NodeSyncer) {
	logger.Info("start node syncing")

	for range state.BtSignal {
		if err := broadcastBt(state, syncer); err != nil {
			logger.Error("node syncing", "error", err)
		}
	}
}

// broadcastBt pull bt and broadcast job to sessions
func broadcastBt(state *ServerState, syncer NodeSyncer) error {
	template, err := syncer.Pull()
	if err != nil {
		return err
	}
	if template != nil && state.updateBlockTemplate(template) {
		// broadcast new job immediately
		state.broadcast()
	}
	return nil
}
