package solana

import (
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const (
	DevnetRPCEndpoint  = rpc.DevnetRPCEndpoint
	TestnetRPCEndpoint = rpc.TestnetRPCEndpoint
	MainnetRPCEndpoint = rpc.MainnetRPCEndpoint
)

type Chain struct {
	*Util
	RpcUrl string

	cli *client.Client
}

func NewChainWithRpc(rpcUrl string) *Chain {
	util := NewUtil()
	return &Chain{
		Util:   util,
		RpcUrl: rpcUrl,
	}
}

func (c *Chain) Client() *client.Client {
	if c.cli == nil {
		c.cli = client.NewClient(c.RpcUrl)
	}
	return c.cli
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() basic.Token {
	return &Token{chain: c}
}

func (c *Chain) BalanceOfAddress(address string) (*basic.Balance, error) {
	client := c.Client()
	balance, err := client.GetBalance(context.Background(), address)
	if err != nil {
		return nil, err
	}
	b := strconv.FormatUint(balance, 10)
	return &basic.Balance{Total: b, TotalWithDecimal: b}, nil
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*basic.Balance, error) {
	address, err := EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return c.BalanceOfAddress(address)
}
func (c *Chain) BalanceOfAccount(account basic.Account) (*basic.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	client := c.Client()
	bytes, err := hexTypes.HexDecodeString(signedTx)
	if err != nil {
		return "", err
	}
	transaction, err := types.TransactionDeserialize(bytes)
	if err != nil {
		return "", err
	}
	res, err := client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Chain) SendSignedTransaction(signedTxn basic.SignedTransaction) (hash *basic.OptionalString, err error) {
	defer basic.CatchPanicAndMapToBasicError(&err)
	txn := AsSignedTransaction(signedTxn)
	if txn == nil {
		return nil, basic.ErrInvalidTransactionType
	}

	client := c.Client()
	res, err := client.SendTransaction(context.Background(), txn.Transaction)
	if err != nil {
		return nil, err
	}
	return &basic.OptionalString{Value: res}, nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (*basic.TransactionDetail, error) {
	response, err := c.Client().GetTransaction(context.Background(), hash)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("not found")
	}
	detail := &basic.TransactionDetail{HashString: hash}
	decodeTransaction(response, detail)
	return detail, nil
}

func (c *Chain) FetchTransactionStatus(hash string) basic.TransactionStatus {
	response, err := c.Client().GetTransaction(context.Background(), hash)
	if err != nil {
		return basic.TransactionStatusNone
	}
	if response == nil || response.Meta == nil {
		return basic.TransactionStatusPending
	}
	if response.Meta.Err == nil {
		return basic.TransactionStatusSuccess
	}
	return basic.TransactionStatusFailure
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := basic.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

func (c *Chain) EstimateTransactionFee(transaction basic.Transaction) (fee *basic.OptionalString, err error) {
	txn, ok := transaction.(*Transaction)
	if !ok {
		return nil, basic.ErrInvalidTransactionType
	}

	client := c.Client()
	gasFee, err := client.GetFeeForMessage(context.Background(), txn.Message)
	if err != nil {
		return nil, err
	}
	feeString := strconv.FormatUint(*gasFee, 10)

	return &basic.OptionalString{Value: feeString}, nil
}
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction basic.Transaction, pubkey string) (fee *basic.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}

func decodeTransaction(tx *client.Transaction, to *basic.TransactionDetail) {
	basic.CatchPanicAndMapToBasicError(nil)

	if tx == nil || tx.BlockTime == nil {
		to.Status = basic.TransactionStatusPending
		return
	}

	to.EstimateFees = strconv.FormatUint(tx.Meta.Fee, 10)
	to.FinishTimestamp = *tx.BlockTime
	if tx.Meta.Err == nil {
		to.Status = basic.TransactionStatusSuccess
	} else {
		to.Status = basic.TransactionStatusFailure
		to.FailureMessage = tx.Meta.LogMessages[len(tx.Meta.LogMessages)-1]
	}

	message := tx.Transaction.Message
	for _, instruction := range message.Instructions {
		program := message.Accounts[instruction.ProgramIDIndex]
		if program != common.SystemProgramID {
			continue
		}

		// We only support decode amount transfer currently.
		data := instruction.Data
		instruct := binary.LittleEndian.Uint32(data[:4])
		toidx := -1
		switch system.Instruction(instruct) {
		case system.InstructionTransfer:
			toidx = instruction.Accounts[1]
		case system.InstructionTransferWithSeed:
			toidx = instruction.Accounts[2]
		default:
			continue
		}

		fromidx := instruction.Accounts[0]
		to.FromAddress = message.Accounts[fromidx].ToBase58()
		to.ToAddress = message.Accounts[toidx].ToBase58()
		amount := binary.LittleEndian.Uint64(data[4:12])
		to.Amount = strconv.FormatUint(amount, 10)
		return
	}
	return
}
