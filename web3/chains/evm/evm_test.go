package evm

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var (
	TestOwnerAddress       = common.HexToAddress("0x3D7FCa6974657CE2658cb5f25D4B091603CfA7e1")
	TestOwnerAddress25     = common.HexToAddress("0x5a127d69aa119e9628731d6caed3350b7c057313")
	TestOwnerPrivateKeyStr = ""
)

func TestNew(t *testing.T) {
	h, err := NewXterioHandler("", TestOwnerPrivateKeyStr, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(h.ChainName(), h.Address.String())

	valueStr := "0.00203"
	for i := 0; i < 3; i++ {
		t := rand.Intn(5) + 1
		valueStr += fmt.Sprintf("%d", t)
	}
	fmt.Println(valueStr)

	bigInt, err := ConvertStringToBigInt(valueStr, 18)
	fmt.Println(bigInt, err)
}

func TestGetEthBalance(t *testing.T) {
	h, err := NewXterioHandler("", TestOwnerPrivateKeyStr, "")
	if err != nil {
		panic(err)
	}
	balance, err := h.GetEthBalance(TestOwnerAddress)
	if err != nil {
		panic(err)
	}
	toString := ConvertBigIntToString(balance, 18)
	fmt.Println(balance.String(), toString)
}

func TestSendEth(t *testing.T) {
	h, err := NewXterioHandler("", TestOwnerPrivateKeyStr, "")
	if err != nil {
		panic(err)
	}
	balance, err := h.GetEthBalance(TestOwnerAddress)
	if err != nil {
		panic(err)
	}
	toString := ConvertBigIntToString(balance, 18)
	fmt.Println(balance.String(), toString)
	amount, err := ConvertStringToBigInt("0.0001", 18)
	if err != nil {
		panic(err)
	}

	tx, err := h.SendEth(TestOwnerAddress25, amount, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.String())
}

func TestCallContract(t *testing.T) {
	h, err := NewXterioHandler("", TestOwnerPrivateKeyStr, "")
	if err != nil {
		panic(err)
	}
	contractAddress := common.HexToAddress("0xBeEDBF1d1908174b4Fc4157aCb128dA4FFa80942")
	abiStr := `[{
        "inputs": [
            {
                "internalType": "uint8",
                "name": "utilityType",
                "type": "uint8"
            }
        ],
        "name": "claimUtility",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    }]`
	method := "claimUtility"

	tx, err := h.CallContract(contractAddress, abiStr, method, nil, big.NewInt(0), uint8(1))
	if err != nil {
		panic(err)
	}
	fmt.Println(tx.String())
}

func TestCallTransaction(t *testing.T) {
	h, err := NewXterioHandler("", TestOwnerPrivateKeyStr, "")
	if err != nil {
		panic(err)
	}
	hash := common.HexToHash("0xcd88a819a43d6b73d35bf1520465fa5040b4de81ea135ad5a0f7272ebe637de8")

	pending, status, err := h.QueryTransactionStatus(hash)
	fmt.Println(pending, status, err)
}
