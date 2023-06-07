package myWeb3

import (
	"context"
	"crypto/ecdsa"
	"errors"

	ERC20 "goTool/myWeb3/erc20"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetTokenBalance(chain string, address string, token string) (res string, err error) {
	var client *ethclient.Client

	if client, err = GetRpcClientByName(chain); err != nil {
		return "", err
	}

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
