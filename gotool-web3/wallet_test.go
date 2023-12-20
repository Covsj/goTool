package gotool_web3

import (
	"fmt"
	"testing"
)

func TestSignText(t *testing.T) {
	//0x3ed07d9ee04f383ea09ed024df5c94386981db23ea9a0cf80eb0482b6d588eb14fa25991726d02bf08950c443786b2f7b9c657c4f426ba04123c2805fb44199b1b
	//0x3ed07d9ee04f383ea09ed024df5c94386981db23ea9a0cf80eb0482b6d588eb14fa25991726d02bf08950c443786b2f7b9c657c4f426ba04123c2805fb44199b00

	key, address, err := CreateNewWallet()
	if err != nil {
		panic(err)
	}
	fmt.Println(key, address)
	sign, err := SignText(key, "skygate")
	if err != nil {
		panic(err)
	}
	fmt.Println(sign)

}
