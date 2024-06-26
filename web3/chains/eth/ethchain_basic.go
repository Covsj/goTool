package eth

import (
	"context"
	"strconv"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/ethereum/go-ethereum/common"
)

// Balance 主网代币余额查询
func (e *EthChain) Balance(address string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	result, err := e.RemoteRpcClient.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return "0", basic.MapAnyToBasicError(err)
	}
	return result.String(), nil
}

// LatestBlockNumber 获取最新区块高度
func (e *EthChain) LatestBlockNumber() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	number, err := e.RemoteRpcClient.BlockNumber(ctx)
	if err != nil {
		return 0, basic.MapAnyToBasicError(err)
	}

	return int64(number), nil
}

// Nonce 获取账户nonce
func (e *EthChain) Nonce(spenderAddressHex string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	nonce, err := e.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return "0", basic.MapAnyToBasicError(err)
	}
	return strconv.FormatUint(nonce, 10), nil
}

// GetChainId 获取链ID
func GetChainId(e *EthChain) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	chainId, err := e.RemoteRpcClient.ChainID(ctx)
	if err != nil {
		return "0", basic.MapAnyToBasicError(err)
	}

	return chainId.String(), nil
}

// SendRawTransaction 对交易进行广播
func (e *EthChain) SendRawTransaction(txHex string) (string, error) {
	var hash common.Hash
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	err := e.RemoteRpcClient.Client().CallContext(ctx, &hash, "eth_sendRawTransaction", txHex)
	if err != nil {
		return "", basic.MapAnyToBasicError(err)
	}
	return hash.String(), nil
}
