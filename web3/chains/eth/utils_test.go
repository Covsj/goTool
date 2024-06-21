package eth

import (
	"fmt"
	"testing"
)

func TestCreateNewWalletByMnemonic(t *testing.T) {
	mnemonic, key, address, err := CreateNewWalletByMnemonic(0)
	fmt.Println(mnemonic)
	fmt.Println(key)
	fmt.Println(address)
	fmt.Println(err)
}

func TestConvertStringToBigInt(t *testing.T) {
	resBigInt, err := BalanceStringToBigInt("0.1", 18)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resBigInt.String())

	toString := BalanceBigIntToString(resBigInt, 18)
	fmt.Println(toString)
}
