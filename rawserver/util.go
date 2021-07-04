package rawserver

import (
	"encoding/hex"
	"math/big"
)

var Diff1 = StringToBig("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

func StringToBig(h string) *big.Int {
	n := new(big.Int)
	n.SetString(h, 0)
	return n
}

func GetTargetHex(diff int64) string {
	padded := make([]byte, 32)

	diffBuff := new(big.Int).Div(Diff1, big.NewInt(diff)).Bytes()
	copy(padded[32-len(diffBuff):], diffBuff)
	buff := padded[0:4]
	targetHex := hex.EncodeToString(reverse(buff))
	return targetHex
}

func reverse(src []byte) []byte {
	dst := make([]byte, len(src))
	for i := len(src); i > 0; i-- {
		dst[len(src)-i] = src[i-1]
	}
	return dst
}
