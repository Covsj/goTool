package basic

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreateNewWalletByMnemonic(t *testing.T) {
	mnemonic, key, address, err := CreateNewWalletByMnemonic()
	fmt.Println(mnemonic)
	fmt.Println(key)
	fmt.Println(address)
	fmt.Println(err)
}

func TestCalculateLastWord(t *testing.T) {
	str := `trigger trigger skill toast machine wave pact black remain poverty slot`
	word, err := CalculateLastWord(strings.Split(str, " ")[:11])
	fmt.Println(word, err)
}
