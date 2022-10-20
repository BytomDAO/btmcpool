package superstratum

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/bytom/btmcpool/common/format/generated"
	"github.com/bytom/btmcpool/common/logger"
)

func TestBanMiner(t *testing.T) {
	t.Skip("Skipping testing in CI environment")

	err := logger.Init(logger.DebugLevel)
	assert.NoError(t, err)

	connCtl := NewConnCtl(10*time.Second, pb.CoinType_BTMC, false, 2, 2, 1.1, 1.1, nil)
	miner := "xx.yy"

	_, hit := connCtl.QueryBanMiner(miner)
	assert.False(t, hit)

	connCtl.AddBanMiner(miner, connCtl.defaultBanPeriod)
	timeout, hit := connCtl.QueryBanMiner(miner)
	assert.True(t, hit)
	assert.True(t, time.Now().Before(timeout.(time.Time)))

	connCtl.RemoveBanMiner(miner)
	_, hit = connCtl.QueryBanMiner(miner)
	assert.False(t, hit)

	// ========== session ctl ==========
	sessionCtl := NewSessionCtl()
	sessionCtl.MinerErrCnt = 100
	sessionCtl.MinerAcCnt = 10
	err = connCtl.JudgeMiner(miner, sessionCtl)
	assert.Error(t, err)

	timeout, hit = connCtl.QueryBanMiner(miner)
	assert.True(t, hit)
	assert.True(t, time.Now().Before(timeout.(time.Time)))

	time.Sleep(10 * time.Second)

	err = connCtl.JudgeMiner(miner, sessionCtl)
	assert.NoError(t, err)
	_, hit = connCtl.QueryBanMiner(miner)
	assert.False(t, hit)
}
