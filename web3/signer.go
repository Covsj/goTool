package web3

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func SignerTransaction(priKey, fromAddress, chainId string, to, gasPrice, gasLimit, value, data string) (string, error) {
	client, err := GetRpcClientByChainId(chainId)
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(priKey)
	if err != nil {
		return "", err
	}
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		return "", err
	}
	v, _ := new(big.Int).SetString(value, 10)
	limit, _ := strconv.ParseUint(gasLimit, 10, 64)
	price, _ := new(big.Int).SetString(gasPrice, 10)

	tx := types.NewTransaction(nonce, common.HexToAddress(to),
		v, limit, price, []byte(data))
	chain, _ := new(big.Int).SetString(chainId, 10)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chain), privateKey)
	if err != nil {
		return "", nil
	}

	b, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil

}
