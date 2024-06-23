package basic

import (
	"strconv"
)

type Balance struct {
	Total            string
	TotalWithDecimal string
}

func EmptyBalance() *Balance {
	return &Balance{
		Total:            "0",
		TotalWithDecimal: "0",
	}
}

func NewBalance(amount string) *Balance {
	return &Balance{Total: amount, TotalWithDecimal: amount}
}

func NewBalanceWithInt(amount int64) *Balance {
	a := strconv.FormatInt(amount, 10)
	return &Balance{Total: a, TotalWithDecimal: a}
}
