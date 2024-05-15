package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Handler interface {
	ChainId() int64
	DefaultGasLimit() uint64
	ChainName() string
	Client() *ethclient.Client
	RpcClient() *rpc.Client
	GetEthBalance(address common.Address) (*big.Int, error)
	GetTokenBalance(tokenAddress common.Address, ownerAddress common.Address) (*big.Int, error)
	SendTransaction(toAddress common.Address, amount *big.Int, gasLimit uint64, data []byte, result interface{}) (common.Hash, error)
	SendEth(toAddress common.Address, amount *big.Int, gasLimit uint64) (common.Hash, error)
	SendToken(tokenAddress common.Address, toAddress common.Address, amount *big.Int) (common.Hash, error)
	GetTokenDecimals(tokenAddress common.Address) (uint8, error)
	CallContract(contractAddress common.Address, abiStr, method string, results interface{},
		amount *big.Int, params ...interface{}) (common.Hash, error)
	QueryTransactionStatus(txHash common.Hash) (isPending bool, status bool, err error)
}
