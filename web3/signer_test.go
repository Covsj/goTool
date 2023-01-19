package web3

import (
	"fmt"
	"testing"
)

func TestSignerTransaction(t *testing.T) {
	prikey, address, _ := GetNewWallet()
	res, err := SignerTransaction(prikey,
		address, "1", "0x15874d65e649880c2614e7a480cb7c9a55787ff6", "18414010547", "200000", "0",
		"0x095ea7b30000000000000000000000001111111254eeb25477b68fb85ed929f73a960582ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	fmt.Println(res, err)
}
