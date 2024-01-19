package gotool_web3

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
