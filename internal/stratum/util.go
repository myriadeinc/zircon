package stratum

import (
	"encoding/hex"
	"math/big"
	"strconv"
)

var maxValue = stringToBig("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

const difficulty500 = "6e128300"

// Although we accept a string, we expect it to be an integer
func convertDifficultyToHex(difficulty string) string {
	diff, err := strconv.ParseInt(difficulty, 10, 64)
	if err != nil {
		return difficulty500
	}

	padded := make([]byte, 32)
	diffBuff := new(big.Int).Div(maxValue, big.NewInt(diff)).Bytes()
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

func stringToBig(h string) *big.Int {
	n := new(big.Int)
	n.SetString(h, 0)
	return n
}
