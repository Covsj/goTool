package basic

import (
	"fmt"
	"testing"
)

func TestCreateNewWalletByMnemonic(t *testing.T) {
	mnemonic, key, address, err := CreateNewWalletByMnemonic()
	fmt.Println(mnemonic)
	fmt.Println(key)
	fmt.Println(address)
	fmt.Println(err)
}
