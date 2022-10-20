package util

import (
	"encoding/hex"
	"hash/crc32"
	"math/big"
)

const dividend = "0x1FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"

func HashMiner(x string) string {
	v := int(crc32.ChecksumIEEE([]byte(x)))
	if -v >= 0 {
		v = -v
	}

	// cut the length of miner to 12
	ans := ""
	for i := 0; i < 12; i++ {
		tmp := v%10 + 97 // ascii (97 -- a)
		v /= 10
		ans += string(tmp)
	}
	return ans
}

func GetTargetHex(diff *big.Int) string {
	padded := make([]byte, 32)
	diffBuff := new(big.Int).Div(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil), diff).Bytes()
	copy(padded[32-len(diffBuff):], diffBuff)
	buff := padded[0:4]
	targetHex := hex.EncodeToString(reverse(buff))
	return targetHex
}

func GetDividend() *big.Int {
	return stringToBig(dividend)
}

func stringToBig(h string) *big.Int {
	n := new(big.Int)
	n.SetString(h, 0)
	return n
}

func reverse(src []byte) []byte {
	dst := make([]byte, len(src))
	for i := len(src); i > 0; i-- {
		dst[len(src)-i] = src[i-1]
	}
	return dst
}
