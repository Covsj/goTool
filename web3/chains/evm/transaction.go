package evm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func CreateTransaction(h Handler, address string, to, gasLimit, gasPrice, value, data string) (*types.Transaction, error) {
	nonce, err := WalletNonce(h, address)
	if err != nil {
		return nil, err
	}
	//    gasPrice := "20000000000" // 例如：20 Gwei
	//    gasLimit := "21000"       // 例如：21000
	//    value := "1000000000000000000" // 例如：1 ETH
	if gasLimit == "" {
		gasLimit = "21000"
	}
	if gasPrice == "" {
		gasPrice = "20000000000"
	}
	valueBigInt, _ := new(big.Int).SetString(value, 10)
	gasLimitUint64, _ := strconv.ParseUint(gasLimit, 10, 64)
	gasPriceBigInt, _ := new(big.Int).SetString(gasPrice, 10)

	tx := types.NewTransaction(nonce, common.HexToAddress(to), valueBigInt, gasLimitUint64, gasPriceBigInt, []byte(data))

	return tx, nil
}

func SignTransactionData(privateKeyHex string, tx *types.Transaction, chainID string) (string, error) {
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}
	// Decode the hex string to a byte slice
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Convert to ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("cannot parse private key: %v", err)
	}

	chainIDBigInt, _ := new(big.Int).SetString(chainID, 10)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainIDBigInt), privateKey)
	if err != nil {
		return "", err
	}

	rlpEncodedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rlpEncodedTx), nil
}
