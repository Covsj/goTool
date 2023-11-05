package gotool_web3

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"

	ERC20 "github.com/Covsj/goTool/gotool-lib/erc20"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetTokenBalance(chain string, address string, token string) (res string, err error) {
	var client *ethclient.Client

	client = GetRpcClientByName(chain)

	// 获取原生token
	if token == "" {
		balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			return "", err
		}
		return balance.String(), nil
	} else {
		tokenAddress := common.HexToAddress(token)
		instance, err := ERC20.NewUtils(tokenAddress, client)
		if err != nil {
			return "", err
		}
		addr := common.HexToAddress(address)

		balance, err := instance.BalanceOf(&bind.CallOpts{}, addr)
		if err != nil {
			return "", err
		}
		return balance.String(), nil
	}
}

func CreateNewWallet() (priKey string, address string, err error) {
	privateKey, err := crypto.GenerateKey()

	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priKey = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return priKey, address, nil
}

func SignerTransaction(priKey, fromAddress, chainId string, to, gasPrice, gasLimit, value, data string) (string, error) {
	client := GetRpcClientByChainId(chainId)
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
