package myUtils

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	RpcEth    = "https://rpc.flashbots.net"
	RpcBsc    = "https://rpc.ankr.com/bsc"
	RPcArb    = "https://arb-mainnet.g.alchemy.com/v2/ezrVAEkTBEkSaXmW_eidWd3ZwrIG0Hfv"
	RpcMatic  = "https://polygon.llamarpc.com"
	RpcZkSync = "https://mainnet.era.zksync.io"
)

var (
	BscClient    *ethclient.Client
	EthClient    *ethclient.Client
	ArbClient    *ethclient.Client
	MaticClient  *ethclient.Client
	ZkSyncClient *ethclient.Client
)

func init() {
	_ = NewWebClient(true)
}

func NewWebClient(needProcess bool) error {
	client, err := NewRpcClient(RpcEth)
	if err != nil {
		if needProcess {
			panic(err)
		} else {
			return err
		}
	}
	EthClient = client

	client, err = NewRpcClient(RpcBsc)
	if err != nil {
		if needProcess {
			panic(err)
		} else {
			return err
		}
	}
	BscClient = client

	client, err = NewRpcClient(RPcArb)
	if err != nil {
		if needProcess {
			panic(err)
		} else {
			return err
		}
	}
	ArbClient = client

	client, err = NewRpcClient(RpcMatic)
	if err != nil {
		if needProcess {
			panic(err)
		} else {
			return err
		}
	}
	MaticClient = client

	client, err = NewRpcClient(RpcZkSync)
	if err != nil {
		if needProcess {
			panic(err)
		} else {
			return err
		}
	}
	ZkSyncClient = client

	return nil
}

func NewRpcClient(rpc string) (client *ethclient.Client, err error) {
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

func GetRpcClientByChainId(chainId string) (client *ethclient.Client) {
	switch chainId {
	case "56":
		client = BscClient
	case "137":
		client = MaticClient
	default:
		client = EthClient
	}
	return client
}

func GetRpcClientByName(chain string) (client *ethclient.Client) {
	switch chain {
	case "bsc":
		client = BscClient
	case "matic":
		client = MaticClient
	case "arb":
		client = ArbClient
	default:
		client = EthClient
	}
	return client
}
