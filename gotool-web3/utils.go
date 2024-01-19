package gotool_web3

import (
	"fmt"
	"math/big"
)

func findIndex(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func binaryToDecimal(binaryStr string) (int, error) {
	num := new(big.Int)
	num, ok := num.SetString(binaryStr, 2)
	if !ok {
		return 0, fmt.Errorf("invalid binary number: %s", binaryStr)
	}
	return int(num.Int64()), nil
}

func binaryStringToBytes(s string) []byte {
	bi := new(big.Int)
	bi.SetString(s, 2)
	return bi.Bytes()
}
