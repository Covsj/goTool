package solana

import (
	"encoding/hex"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
)

type Transaction struct {
	Message types.Message
}

func NewTransaction(rawBase58String string) (*Transaction, error) {
	data, err := base58.Decode(rawBase58String)
	if err != nil {
		return nil, err
	}
	msg, err := types.MessageDeserialize(data)
	if err != nil {
		return nil, err
	}
	return &Transaction{Message: msg}, nil
}

func (t *Transaction) SignWithAccount(account basic.Account) (signedTxn *basic.OptionalString, err error) {
	txn, err := t.SignedTransactionWithAccount(account)
	if err != nil {
		return nil, err
	}
	return txn.HexString()
}

func (t *Transaction) SignedTransactionWithAccount(account basic.Account) (signedTxn basic.SignedTransaction, err error) {
	solanaAcc := AsSolanaAccount(account)
	if solanaAcc == nil {
		return nil, basic.ErrInvalidAccountType
	}

	// create tx by message + signer
	txn, err := types.NewTransaction(types.NewTransactionParam{
		Message: t.Message,
		Signers: []types.Account{*solanaAcc.account},
	})
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{
		Transaction: txn,
	}, nil
}

type SignedTransaction struct {
	Transaction types.Transaction
}

func (txn *SignedTransaction) HexString() (res *basic.OptionalString, err error) {
	txnBytes, err := txn.Transaction.Serialize()
	if err != nil {
		return nil, err
	}
	hexString := "0x" + hex.EncodeToString(txnBytes)

	return &basic.OptionalString{Value: hexString}, nil
}

func AsSignedTransaction(txn basic.SignedTransaction) *SignedTransaction {
	if res, ok := txn.(*SignedTransaction); ok {
		return res
	}
	return nil
}
