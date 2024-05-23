package eth

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c *Chain) ChainId() (string, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "0", err
	}
	return chain.chainId.String(), nil
}

// 获取交易的 nonce
func (c *Chain) NonceOfAddress(address string) (string, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.Nonce(address)
}

func (c *Chain) LatestBlockNumber() (s string, err error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return
	}
	n, err := chain.LatestBlockNumber()
	if err != nil {
		return
	}
	return strconv.FormatInt(n, 10), nil
}

// BatchFetchErc20TokenBalance SDK 批量请求代币余额
func (c *Chain) BatchFetchErc20TokenBalance(contractListString []string, address string) (string, error) {
	//contractList := strings.Split(contractListString, ",")
	balances, err := c.BatchErc20TokenBalance(contractListString, address)
	if err != nil {
		return "", err
	}
	return strings.Join(balances, ","), nil
}

// 批量请求代币余额
// @param contractList 批量查询的代币的合约地址数组
// @param address 用户的钱包地址
// @return 余额数组，顺序与传入的 contractList 是保持一致的
// @throw 如果任意一个代币请求余额出错时，会抛出错误
func (c *Chain) BatchErc20TokenBalance(contractList []string, address string) ([]string, error) {
	return basic.MapListConcurrentStringToString(contractList, func(s string) (string, error) {
		b, err := c.Erc20Token(s).BalanceOfAddress(address)
		return b.Total, err
	})
}

// call eth_call method
// @param blockNumber Especially -2 is the latest block, -1 is pending block.
func (c *Chain) CallContract(msg *CallMsg, blockNumber int64) (string, error) {
	if msg == nil {
		return "", errors.New("Invalid call message")
	}
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), chain.timeout)
	defer cancel()

	var block *big.Int = nil
	if blockNumber >= -1 {
		block = new(big.Int).SetInt64(blockNumber)
	}
	hash, err := chain.RemoteRpcClient.CallContract(ctx, msg.msg, block)
	if err != nil {
		return "", err
	}

	return types.HexEncodeToString(hash), nil
}

func (c *Chain) PendingCallContract(msg *CallMsg) (string, error) {
	return c.CallContract(msg, -1)
}

func (c *Chain) LatestCallContract(msg *CallMsg) (string, error) {
	return c.CallContract(msg, -2)
}

// Sign a transaction
// @return signed tx hash
func (c *Chain) SignTransaction(privateKey string, transaction *Transaction) (*basic.OptionalString, error) {
	privatekeyData, err := types.HexDecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return c.SignWithPrivateKeyData(privatekeyData, transaction)
}

func (c *Chain) SignTransactionWithAccount(account basic.Account, transaction *Transaction) (*basic.OptionalString, error) {
	if account == nil {
		return nil, errors.New("Invalid account")
	}
	privatekeyData, err := account.PrivateKey()
	if err != nil {
		return nil, err
	}
	return c.SignWithPrivateKeyData(privatekeyData, transaction)
}

func (c *Chain) SignWithPrivateKeyData(privateKeyData []byte, transaction *Transaction) (o *basic.OptionalString, err error) {
	if privateKeyData == nil || transaction == nil {
		return nil, errors.New("Invalid privatekey or transaction")
	}

	client, err := GetConnection(c.RpcUrl)
	if err != nil {
		return
	}

	privateKeyECDSA, err := crypto.ToECDSA(privateKeyData)
	if err != nil {
		return
	}

	rawTx, err := transaction.GetRawTx()
	if err != nil {
		return
	}

	output, err := client.buildTxWithTransaction(rawTx, privateKeyECDSA)
	if err != nil {
		return
	}

	return &basic.OptionalString{Value: output.TxHex}, nil
}
