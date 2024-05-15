package ton

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xssnick/tonutils-go/tlb"
)

func TestTonChainHandler_CreateNewWallet(t *testing.T) {
	p, err := NewTonChainHandler()
	if err != nil {
		panic(err)
	}
	wallet, words, err := p.CreateNewWallet([]string{})
	if err != nil {
		panic(err)
	}
	fmt.Println(words)
	fmt.Println(wallet.WalletAddress().String())
}

func TestGetTonBalance(t *testing.T) {
	//log.SetTextFormatter()
	p, err := NewTonChainHandler()
	if err != nil {
		panic(err)
	}
	balance, err := p.GetTonBalance("UQB_2hgiLMfGAiDHTkFZ1CrYfa1AGNoiDsH-lkDw4lN70bFy")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance.String())
}

func TestTransfer(t *testing.T) {
	p, err := NewTonChainHandler()
	if err != nil {
		panic(err)
	}
	w, _, err := p.CreateNewWallet(strings.Split("", " "))
	if err != nil {
		panic(err)
	}
	err = p.Transfer(w, tlb.MustFromTON("0.000000002"), "")
	if err != nil {
		panic(err)
	}
}
