package web3

import (
	"fmt"
	"testing"
)

func TestGetTokenBalance(t *testing.T) {
	addr := "0x67197a52f23236c9a123c266cd6ae114dcc56f90"

	fmt.Println(GetTokenBalance("bsc", addr, ""))
	fmt.Println(GetTokenBalance("bsc", addr, "0x55d398326f99059ff775485246999027b3197955"))
	fmt.Println(GetTokenBalance("", addr, ""))
	fmt.Println(GetTokenBalance("eth", addr, ""))
	fmt.Println(GetTokenBalance("matic", addr, ""))
}

func TestGetNewWallet(t *testing.T) {
	fmt.Println(GetNewWallet())
}
