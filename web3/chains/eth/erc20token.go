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

type Erc20TokenInfo struct {
	*basic.TokenInfo

	ContractAddress string
	ChainId         string
	TokenIcon       string

	Balance string
}

type Erc20Token struct {
	*Token
	ContractAddress string
}

func NewErc20Token(chain *Chain, contractAddress string) *Erc20Token {
	return &Erc20Token{
		Token:           &Token{chain: chain},
		ContractAddress: contractAddress,
	}
}

func (t *Erc20Token) Erc20TokenInfo() (*Erc20TokenInfo, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}
	baseInfo, err := t.TokenInfo()
	if err != nil {
		return nil, err
	}

	return &Erc20TokenInfo{
		TokenInfo:       baseInfo,
		ContractAddress: t.ContractAddress,
		ChainId:         chain.chainId.String(),
	}, nil
}

func (t *Erc20Token) TokenInfo() (*basic.TokenInfo, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}

	info := &basic.TokenInfo{}
	info.Name, err = chain.TokenName(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	info.Symbol, err = chain.TokenSymbol(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	info.Decimal, err = chain.TokenDecimal(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (t *Erc20Token) Decimal() (int16, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return 18, err
	}
	return chain.TokenDecimal(t.ContractAddress)
}

func (t *Erc20Token) BalanceOfAddress(address string) (*basic.Balance, error) {
	b := basic.NewBalance("0")
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return b, err
	}

	balance, err := chain.TokenBalance(t.ContractAddress, address)
	if err != nil {
		return b, err
	}
	return &basic.Balance{
		Total:            balance,
		TotalWithDecimal: balance,
	}, nil
}

func (t *Erc20Token) EstimateGasFeeLayer2(msg *CallMsg) (*OptimismLayer2Gas, error) {
	data, err := EncodeErc20Transfer(msg.GetTo(), msg.GetValue())
	if err != nil {
		return nil, err
	}
	msg.SetData(data)
	return t.Token.EstimateGasFeeLayer2(msg)
}

func (t *Erc20Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	msg := NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(t.ContractAddress)
	msg.SetGasPrice(gasPrice)
	msg.SetValue("0")

	data, err := EncodeErc20Transfer(receiverAddress, amount)
	if err != nil {
		return "", err
	}
	msg.SetData(data)

	gasLimit, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		return "", err
	}
	return gasLimit.Value, nil
}

func (t *Erc20Token) BuildTransferTx(privateKey string, transaction *Transaction) (*basic.OptionalString, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return t.buildTransfer(privateKeyECDSA, transaction)
}

func (t *Erc20Token) Transfer(account *Account, toAddress, amount string) (string, error) {
	err := errors.New("transfer failed")
	transferData, err := EncodeErc20Transfer(toAddress, amount)
	if err != nil {
		return "", err
	}
	gasPrice, err := t.chain.SuggestGasPrice()
	if err != nil {
		return "", err
	}
	msg := NewCallMsg()
	msg.SetFrom(account.Address())
	msg.SetTo(t.ContractAddress)
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(transferData)
	msg.SetValue("0")
	gasLimit, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		gasLimit = &basic.OptionalString{Value: "100000"}
		err = nil
	}
	msg.SetGasLimit(gasLimit.Value)
	transaction := msg.TransferToTransaction()
	rawTx, err := t.chain.buildTransfer(account.privateKeyECDSA, transaction)
	if err != nil {
		return "", err
	}
	return t.chain.SendRawTransaction(rawTx.Value)
}

func (t *Erc20Token) BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*basic.OptionalString, error) {
	return t.buildTransfer(account.privateKeyECDSA, transaction)
}

func (t *Erc20Token) buildTransfer(privateKey *ecdsa.PrivateKey, transaction *Transaction) (*basic.OptionalString, error) {
	err := transaction.TransformToErc20Transaction(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	return t.chain.buildTransfer(privateKey, transaction)
}

func (t *Erc20Token) Allowance(owner, spender string) (*big.Int, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}
	res := big.NewInt(0)
	ownerAddress := common.HexToAddress(owner)
	spenderAddress := common.HexToAddress(spender)
	err = chain.CallContractConstant(&res, t.ContractAddress, Erc20AbiStr, "allowance", nil, ownerAddress, spenderAddress)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *Erc20Token) Approve(account *Account, spender string, amount *big.Int) (string, error) {
	err := errors.New("approve failed")

	approveData, err := EncodeErc20Approve(spender, amount)
	if err != nil {
		return "", err
	}

	gasPrice, err := t.chain.SuggestGasPrice()
	if err != nil {
		return "", err
	}
	msg := NewCallMsg()
	msg.SetFrom(account.Address())
	msg.SetTo(t.ContractAddress)
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(approveData)
	msg.SetValue("0")

	gasLimit, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		gasLimit = &basic.OptionalString{Value: "100000"}
		err = nil
	}
	msg.SetGasLimit(gasLimit.Value)

	transaction := msg.TransferToTransaction()
	rawTx, err := t.chain.buildTransfer(account.privateKeyECDSA, transaction)
	if err != nil {
		return "", err
	}

	return t.chain.SendRawTransaction(rawTx.Value)
}

func (t *Erc20Token) BuildTransfer(sender, receiver, amount string) (txn basic.Transaction, err error) {
	return nil, basic.ErrUnsupportedFunction
}
func (t *Erc20Token) CanTransferAll() bool {
	return false
}
func (t *Erc20Token) BuildTransferAll(sender, receiver string) (txn basic.Transaction, err error) {
	return nil, basic.ErrUnsupportedFunction
}
