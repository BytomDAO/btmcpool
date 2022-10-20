package superstratum

import (
	"sync/atomic"
	"time"

	"github.com/cornelk/hashmap"
	"golang.org/x/time/rate"

	pb "github.com/bytom/btmcpool/common/format/generated"
	"github.com/bytom/btmcpool/common/logger"
)

// connection control, manage the blacklist and accessing privilege

type connCtrl struct {
	ipBanEnable      bool
	ipLimiterList    *hashmap.HashMap // ip => limiter, to restrict throughput and connection
	ipBanList        *hashmap.HashMap
	minerBanList     *hashmap.HashMap
	minerBanCnt      int32
	defaultBanPeriod time.Duration
	coin             pb.CoinType
	fwdTrigger       bool
	maxThroughput    int
	maxConnection    int
	burstThroughput  int
	burstConnection  int
	whiteIpMap       map[string]bool
}

func NewConnCtl(banPeriod time.Duration, coin pb.CoinType, ipBanEnable bool, maxThroughput, maxConnection int, throughputRatio, ConnectionRatio float64, whiteList []string) *connCtrl {
	whiteIpMap := make(map[string]bool)
	for _, ip := range whiteList {
		whiteIpMap[ip] = true
	}

	return &connCtrl{
		ipBanEnable:      ipBanEnable,
		ipLimiterList:    &hashmap.HashMap{},
		ipBanList:        &hashmap.HashMap{},
		minerBanList:     &hashmap.HashMap{},
		minerBanCnt:      0,
		defaultBanPeriod: banPeriod,
		coin:             coin,
		maxThroughput:    maxThroughput,
		maxConnection:    maxConnection,
		burstThroughput:  int(float64(maxThroughput) * throughputRatio),
		burstConnection:  int(float64(maxConnection) * ConnectionRatio),
		whiteIpMap:       whiteIpMap,
	}
}

type ipLimiter struct {
	throughput *rate.Limiter // throughput limiter
	connection *rate.Limiter // connection limiter
}

type SessionCtl struct {
	MinerAcCnt  uint64
	MinerErrCnt uint64
}

func NewSessionCtl() *SessionCtl {
	return &SessionCtl{
		MinerAcCnt:  0,
		MinerErrCnt: 0,
	}
}

func (s *connCtrl) QueryBanMiner(miner string) (interface{}, bool) {
	return s.minerBanList.Get(miner)
}

func (s *connCtrl) AddBanMiner(miner string, bannedDur time.Duration) {
	// banPeriod means disable ban mechanism
	if s.defaultBanPeriod.Nanoseconds() == int64(0) {
		return
	}

	logger.Info("add banned miner", "miner", miner, "banned duration", bannedDur)
	s.minerBanList.Set(miner, time.Now().Add(bannedDur))
	atomic.AddInt32(&s.minerBanCnt, 1)
}

func (s *connCtrl) RemoveBanMiner(miner string) {
	logger.Info("remove banned miner", "miner", miner)
	s.minerBanList.Del(miner)
	atomic.AddInt32(&s.minerBanCnt, -1)
}

func (s *connCtrl) JudgeMiner(miner string, c *SessionCtl) error {
	// banPeriod means disable ban mechanism
	if s.defaultBanPeriod.Nanoseconds() == int64(0) {
		return nil
	}

	timeout, hit := s.QueryBanMiner(miner)
	if hit {
		if timeout.(time.Time).After(time.Now()) {
			return ErrBannedMiner
		} else {
			s.RemoveBanMiner(miner)
		}
	} else {
		// this part should judge in coin related code, lazy & hard code here
		if (c.MinerAcCnt+c.MinerErrCnt >= 100) && float64(c.MinerErrCnt)/float64(c.MinerAcCnt+c.MinerErrCnt) >= 0.5 {
			s.AddBanMiner(miner, s.defaultBanPeriod)
			c.MinerErrCnt = 0
			c.MinerAcCnt = 0
			return ErrBannedMiner
		}
	}
	return nil
}

func (s *connCtrl) scanning(period time.Duration) {
	for {
		for m := range s.minerBanList.Iter() {
			if m.Value.(time.Time).Add(30 * time.Minute).Before(time.Now()) {
				// if expire endTime more than 30min, miners are still in banList
				// remove it manually
				s.RemoveBanMiner(m.Key.(string))
			}
		}
		time.Sleep(period)
	}
}
