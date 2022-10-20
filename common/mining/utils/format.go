package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func ToLittleEndianHex(source uint64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, source)
	return fmt.Sprintf("%x", b)
}

func FromLittleEndianHex(s string) (uint64, error) {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(decoded), nil
}
