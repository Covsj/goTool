package basic

import (
	"encoding/json"
	"strings"
)

type Transaction interface {
	SignWithAccount(account Account) (signedTxn *OptionalString, err error)
	SignedTransactionWithAccount(account Account) (signedTxn SignedTransaction, err error)
}

type SignedTransaction interface {
	HexString() (res *OptionalString, err error)
}

type TransactionStatus = SDKEnumInt

const (
	TransactionStatusNone    TransactionStatus = 0
	TransactionStatusPending TransactionStatus = 1
	TransactionStatusSuccess TransactionStatus = 2
	TransactionStatusFailure TransactionStatus = 3
)

type TransactionDetail struct {
	// 哈希字符串
	HashString string

	// 交易数量
	Amount string

	EstimateFees string

	// 发送地址
	FromAddress string
	// 接收地址
	ToAddress string

	// 交易状态
	Status TransactionStatus
	// 完成时间戳，如果是Pending状态则是0
	FinishTimestamp int64
	// 失败信息
	FailureMessage string

	// 如果此交易是 CID 转账，则其值将为 CID，否则为空
	CIDNumber string
	// 如果此交易是 NFT 转账，其值将为 Token 名称，否则为空
	TokenName string
}

func (d *TransactionDetail) IsCIDTransfer() bool {
	return strings.TrimSpace(d.CIDNumber) != ""
}

func (d *TransactionDetail) IsNFTTransfer() bool {
	return strings.TrimSpace(d.TokenName) != ""
}

func (d *TransactionDetail) JsonString() string {
	b, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(b)
}
