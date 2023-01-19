package web3

import (
	"errors"

	"goTool/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(rpc string) (client *ethclient.Client, err error) {
	if rpc == "" {
		return nil, errors.New("rpc is null")
	}
	for i := 0; i < 3; i++ {
		client, err = ethclient.Dial(rpc)
		if err != nil {
			continue
		} else {
			return client, nil
		}
	}
	return nil, err
}

func GetRpcClientByName(chain string) (client *ethclient.Client, err error) {
	switch chain {
	case "bsc":
		if config.DefaultToolConfig.Web3Config.BscClient == nil {
			for i := 0; i < 3; i++ {
				client, err = ethclient.Dial(config.DefaultToolConfig.Web3Config.BscRpc)
				if err != nil {
					continue
				}
			}
			config.DefaultToolConfig.Web3Config.BscClient = client
		}
		client = config.DefaultToolConfig.Web3Config.BscClient
	case "matic":
		if config.DefaultToolConfig.Web3Config.MaticClient == nil {
			for i := 0; i < 3; i++ {
				client, err = ethclient.Dial(config.DefaultToolConfig.Web3Config.MaticRpc)
				if err != nil {
					continue
				}
			}
			config.DefaultToolConfig.Web3Config.MaticClient = client
		}
		client = config.DefaultToolConfig.Web3Config.MaticClient
	default:
		if config.DefaultToolConfig.Web3Config.EthClient == nil {
			for i := 0; i < 3; i++ {
				client, err = ethclient.Dial(config.DefaultToolConfig.Web3Config.EthRpc)
				if err != nil {
					continue
				}
			}
			config.DefaultToolConfig.Web3Config.EthClient = client
		}
		client = config.DefaultToolConfig.Web3Config.EthClient
	}
	return client, err
}

func GetRpcClientByChainId(chainId string) (client *ethclient.Client, err error) {
	switch chainId {
	case "56":
		if config.DefaultToolConfig.Web3Config.BscClient == nil {
			for i := 0; i < 3; i++ {
				client, err = ethclient.Dial(config.DefaultToolConfig.Web3Config.BscRpc)
				if err != nil {
					continue
				}
			}
			config.DefaultToolConfig.Web3Config.BscClient = client
		}
		client = config.DefaultToolConfig.Web3Config.BscClient
	case "137":
		if config.DefaultToolConfig.Web3Config.MaticClient == nil {
			for i := 0; i < 3; i++ {
				client, err = ethclient.Dial(config.DefaultToolConfig.Web3Config.MaticRpc)
				if err != nil {
					continue
				}
			}
			config.DefaultToolConfig.Web3Config.MaticClient = client
		}
		client = config.DefaultToolConfig.Web3Config.MaticClient
	default:
		if config.DefaultToolConfig.Web3Config.EthClient == nil {
			for i := 0; i < 3; i++ {
				client, err = ethclient.Dial(config.DefaultToolConfig.Web3Config.EthRpc)
				if err != nil {
					continue
				}
			}
			config.DefaultToolConfig.Web3Config.EthClient = client
		}
		client = config.DefaultToolConfig.Web3Config.EthClient
	}
	return client, err
}
