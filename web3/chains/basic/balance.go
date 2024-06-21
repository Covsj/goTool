package basic

import (
	"fmt"
	"math/big"
	"strconv"
)

type Balance struct {
	Total  string
	Usable string
}

func EmptyBalance() *Balance {
	return &Balance{
		Total:  "0",
		Usable: "0",
	}
}

func NewBalance(amount string) *Balance {
	return &Balance{Total: amount, Usable: amount}
}

func NewBalanceWithInt(amount int64) *Balance {
	a := strconv.FormatInt(amount, 10)
	return &Balance{Total: a, Usable: a}
}

func ConvertBigIntToString(balance *big.Int, decimals uint8) string {
	if decimals == 0 {
		decimals = 18
	}
	// Convert balance to float
	balanceFloat := new(big.Float).SetInt(balance)

	// Calculate divisor
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Divide balance by divisor
	balanceFloat.Quo(balanceFloat, divisor)

	// Convert float balance to string
	balanceStr := balanceFloat.Text('f', int(decimals))

	return balanceStr
}

func ConvertStringToBigInt(balanceStr string, decimals uint8) (*big.Int, error) {
	if decimals == 0 {
		decimals = 18
	}
	// Convert balance string to big.Float
	balanceFloat, ok := new(big.Float).SetString(balanceStr)
	if !ok {
		return nil, fmt.Errorf("failed to convert balance to big.Float")
	}

	// Multiply balance by 10^decimals
	exp := big.NewInt(int64(decimals))
	multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), exp, nil))
	balanceFloat.Mul(balanceFloat, multiplier)

	// Convert balance to big.Int
	balanceInt := new(big.Int)
	balanceFloat.Int(balanceInt)

	return balanceInt, nil
}
