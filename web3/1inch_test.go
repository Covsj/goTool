package web3

import (
	"fmt"
	"testing"
)

func TestCheckAllowance(t *testing.T) {
	allowance, err := CheckAllowance("1", "",
		"", "")
	fmt.Println(allowance, err)
}
func TestGetApproveTransaction(t *testing.T) {
	allowance, err := GetApproveTransaction("1", "",
		"")
	fmt.Println(allowance, err)
	allowance, err = GetApproveTransaction("1", "",
		"123123")
	fmt.Println(allowance, err)
}
