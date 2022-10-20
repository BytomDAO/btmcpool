package superstratum

import (
	"math/big"
)

func NewDiffAdjust(diff *big.Int) *diffAdjust {
	return &diffAdjust{diff: diff}
}

type diffAdjust struct {
	diff *big.Int
}

func (s *diffAdjust) GetDiff() *big.Int {
	return s.diff
}
