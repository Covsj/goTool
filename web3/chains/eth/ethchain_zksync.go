package eth

import (
	"context"
	"math/big"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

const zksyncChainId = 324
const zksyncChainIdTestnet = 280

type zksyncTransaction struct {
	BlockHash            common.Hash    `json:"blockHash"`
	BlockNumber          *hexutil.Big   `json:"blockNumber"`
	ChainID              hexutil.Big    `json:"chainId"`
	From                 common.Address `json:"from"`
	Gas                  hexutil.Uint64 `json:"gas"`
	GasPrice             hexutil.Big    `json:"gasPrice"`
	Hash                 common.Hash    `json:"hash"`
	Data                 hexutil.Bytes  `json:"input"`
	L1BatchNumber        hexutil.Big    `json:"l1BatchNumber"`
	L1BatchTxIndex       hexutil.Big    `json:"l1BatchTxIndex"`
	MaxFeePerGas         hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas hexutil.Big    `json:"maxPriorityFeePerGas"`
	Nonce                hexutil.Uint64 `json:"nonce"`
	To                   common.Address `json:"to"`
	TransactionIndex     hexutil.Uint   `json:"transactionIndex"`
	TxType               hexutil.Uint64 `json:"type"`
	Value                hexutil.Big    `json:"value"`
	V                    hexutil.Big    `json:"v"`
	R                    hexutil.Big    `json:"r"`
	S                    hexutil.Big    `json:"s"`
}

func (e *EthChain) zksyncFetchTransactionDetail(hashString string) (detail *basic.TransactionDetail, txn *types.Transaction, err error) {
	defer basic.CatchPanicAndMapToBasicError(&err)

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	var tx *zksyncTransaction
	err = e.RemoteRpcClient.Client().CallContext(ctx, &tx, "eth_getTransactionByHash", hashString)
	if err != nil {
		return
	} else if tx == nil {
		return nil, nil, ethereum.NotFound
	}

	gasFeeInt := big.NewInt(0).Mul(tx.GasPrice.ToInt(), big.NewInt(0).SetUint64(uint64(tx.Gas)))
	detail = &basic.TransactionDetail{
		HashString:   hashString,
		FromAddress:  tx.From.String(),
		ToAddress:    tx.To.String(),
		Amount:       tx.Value.ToInt().String(),
		EstimateFees: gasFeeInt.String(),
	}

	if tx.BlockNumber == nil {
		detail.Status = basic.TransactionStatusPending
		return
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	// 当交易没有处于pending状态时，可以查询receipt信息，即交易是否成功, err为nil时，表示查询成功进入if语句赋值
	receipt, err := e.TransactionReceiptByHash(hashString)
	if err != nil {
		return
	}
	blockHeader, err := e.RemoteRpcClient.HeaderByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return
	}

	if receipt.Status == 0 {
		detail.Status = basic.TransactionStatusFailure
		// get error message
		_, err := e.RemoteRpcClient.CallContract(ctx, ethereum.CallMsg{
			From:       tx.From,
			To:         &tx.To,
			Data:       tx.Data,
			Gas:        uint64(tx.Gas),
			GasPrice:   tx.GasPrice.ToInt(),
			GasFeeCap:  tx.MaxFeePerGas.ToInt(),
			GasTipCap:  tx.MaxPriorityFeePerGas.ToInt(),
			Value:      tx.Value.ToInt(),
			AccessList: nil,
		}, receipt.BlockNumber)
		if err != nil {
			detail.FailureMessage = err.Error()
			err = nil
		}

	} else {
		detail.Status = basic.TransactionStatusSuccess
	}

	effectiveGasPrice := tx.GasPrice.ToInt()
	if receipt.EffectiveGasPrice != nil {
		effectiveGasPrice = receipt.EffectiveGasPrice
	}
	gasFeeInt = big.NewInt(0).Mul(effectiveGasPrice, big.NewInt(0).SetUint64(receipt.GasUsed))
	if receipt.L1Fee != nil {
		gasFeeInt = gasFeeInt.Add(gasFeeInt, receipt.L1Fee)
	}
	detail.EstimateFees = gasFeeInt.String()
	detail.FinishTimestamp = int64(blockHeader.Time)

	txn = types.NewTx(&types.DynamicFeeTx{
		Data: tx.Data, // only use data currently.
	})
	return detail, txn, nil
}
