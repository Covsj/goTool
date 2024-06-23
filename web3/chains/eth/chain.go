package eth

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type IChain interface {
	basic.Chain
	SubmitTransactionData(account basic.Account, to string, data []byte, value string) (string, error)
	GetEthChain() (*EthChain, error)
	EstimateGasLimit(msg *CallMsg) (gas *basic.OptionalString, err error)
}

type Chain struct {
	RpcUrl string
}

func NewChainWithRpc(rpcUrl string) *Chain {
	return &Chain{
		RpcUrl: rpcUrl,
	}
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() basic.Token {
	return &Token{chain: c}
}

func (c *Chain) MainEthToken() TokenProtocol {
	return &Token{chain: c}
}

func (c *Chain) Erc20Token(contractAddress string) TokenProtocol {
	return &Erc20Token{
		Token:           &Token{chain: c},
		ContractAddress: contractAddress,
	}
}

func (c *Chain) BalanceOfAddress(address string) (*basic.Balance, error) {
	b := basic.NewBalance("0")

	if !IsValidAddress(address) {
		return b, errors.New("invalid hex address")
	}

	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return b, err
	}
	balance, err := chain.Balance(address)
	if err != nil {
		return b, err
	}
	return &basic.Balance{
		Total:            balance,
		TotalWithDecimal: balance,
	}, nil
}

func (c *Chain) BalanceOfPublicKey(publicKey string) (*basic.Balance, error) {
	return c.BalanceOfAddress(publicKey)
}

func (c *Chain) BalanceOfAccount(account basic.Account) (*basic.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.SendRawTransaction(signedTx)
}

func (c *Chain) SendSignedTransaction(signedTxn basic.SignedTransaction) (hash *basic.OptionalString, err error) {
	return nil, basic.ErrUnsupportedFunction
}

func (c *Chain) FetchTransactionDetail(hash string) (*basic.TransactionDetail, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}
	detail, txn, err := chain.FetchTransactionDetail(hash)
	if err != nil {
		return nil, err
	}
	if data := txn.Data(); len(data) > 0 {
		method, params, err := DecodeContractParams(Erc20AbiStr, data)
		if err == nil && method == Erc20MethodTransfer {
			detail.ToAddress = params[0].(common.Address).String()
			detail.Amount = params[1].(*big.Int).String()
		}
	}
	return detail, nil
}

func (c *Chain) FetchTransactionStatus(hash string) basic.TransactionStatus {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return basic.TransactionStatusNone
	}
	return chain.FetchTransactionStatus(hash)
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return ""
	}
	return chain.SdkBatchTransactionStatus(hashListString)
}

func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction basic.Transaction, pubkey string) (fee *basic.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}

func (c *Chain) EstimateTransactionFee(transaction basic.Transaction) (fee *basic.OptionalString, err error) {
	return nil, basic.ErrUnsupportedFunction
}

func (c *Chain) BuildTransferTx(privateKey string, transaction *Transaction) (*basic.OptionalString, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return c.buildTransfer(privateKeyECDSA, transaction)
}

func (c *Chain) BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*basic.OptionalString, error) {
	return c.buildTransfer(account.privateKeyECDSA, transaction)
}

func (c *Chain) buildTransfer(privateKey *ecdsa.PrivateKey, transaction *Transaction) (*basic.OptionalString, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}

	if transaction.Nonce == "" || transaction.Nonce == "0" {
		address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
		nonce, err := chain.Nonce(address)
		if err != nil {
			nonce = "0"
			err = nil
		}
		transaction.Nonce = nonce
	}

	rawTx, err := transaction.GetRawTx()
	if err != nil {
		return nil, err
	}

	txResult, err := chain.buildTxWithTransaction(rawTx, privateKey)
	if err != nil {
		return nil, err
	}

	return &basic.OptionalString{Value: txResult.TxHex}, nil
}

func (c *Chain) SubmitTransactionData(account basic.Account, to string, data []byte, value string) (string, error) {
	gasPrice, err := c.SuggestGasPrice()
	if err != nil {
		return "", err
	}
	msg := NewCallMsg()
	msg.SetFrom(account.Address())
	msg.SetTo(to)
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(data)
	msg.SetValue(value)

	gasLimit, err := c.EstimateGasLimit(msg)
	if err != nil {
		gasLimit = &basic.OptionalString{Value: "200000"}
		err = nil
	}
	msg.SetGasLimit(gasLimit.Value)
	tx := msg.TransferToTransaction()
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "", err
	}
	nonce, err := chain.Nonce(account.Address())
	if err != nil {
		nonce = "0"
		err = nil
	}
	tx.Nonce = nonce
	privateKeyHex, err := account.PrivateKeyHex()
	if err != nil {
		return "", err
	}
	signedTx, err := c.SignTransaction(privateKeyHex, tx)
	if err != nil {
		return "", err
	}
	return c.SendRawTransaction(signedTx.Value)
}

func (c *Chain) GetEthChain() (*EthChain, error) {
	return GetConnection(c.RpcUrl)
}
