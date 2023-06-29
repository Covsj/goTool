package dex

import (
	"fmt"
	"testing"
)

func TestCheckAllowance(t *testing.T) {
	allowance, err := CheckAllowance("1", "0x15874d65e649880c2614e7a480cb7c9a55787ff6",
		"0x68cb015699cb565bbe4b08402b6886ac26e43f73", "")
	fmt.Println(allowance, err)
}
func TestGetApproveTransaction(t *testing.T) {
	allowance, err := GetApproveTransaction("1", "0x15874d65e649880c2614e7a480cb7c9a55787ff6",
		"")
	fmt.Println(allowance, err)
	allowance, err = GetApproveTransaction("1", "0x15874d65e649880c2614e7a480cb7c9a55787ff6",
		"123123")
	fmt.Println(allowance, err)
}
