package basic

import (
	"strconv"
)

type Balance struct {
	Total   string
	Decimal int16
}

func EmptyBalance() *Balance {
	return &Balance{
		Total:   "0",
		Decimal: 18,
	}
}

func NewBalance(amount string, decimal int16) *Balance {
	return &Balance{Total: amount, Decimal: decimal}
}

func NewBalanceWithInt(amount int64, decimal int16) *Balance {
	a := strconv.FormatInt(amount, 10)
	return &Balance{Total: a, Decimal: decimal}
}
