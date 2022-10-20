package utils

// GetReward calculates the miner reward from share
func GetReward(shareDiff uint64, netDiff uint64, factor uint64, fee uint64) uint64 {
	return uint64(float64(shareDiff) * float64(factor) * (float64(10000) - float64(fee)) / (float64(netDiff) * float64(10000)))
}
