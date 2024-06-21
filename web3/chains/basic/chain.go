package basic

type Chain interface {
	MainToken() Token

	BalanceOfAddress(address string) (*Balance, error)
	BalanceOfPublicKey(publicKey string) (*Balance, error)
	BalanceOfAccount(account Account) (*Balance, error)

	// SendRawTransaction 发送原始交易广播
	// @return the hex hash string
	SendRawTransaction(signedTx string) (string, error)

	// SendSignedTransaction 发送签名交易广播
	// @return the hex hash string
	SendSignedTransaction(signedTxn SignedTransaction) (*OptionalString, error)

	// FetchTransactionDetail 根据hash查询交易详细信息
	FetchTransactionDetail(hash string) (*TransactionDetail, error)

	// FetchTransactionStatus 根据hash查询交易状态
	FetchTransactionStatus(hash string) TransactionStatus

	// BatchFetchTransactionStatus 批量查询交易状态
	// @param  需要查询的hash拼接字符串，比如 "hash1,hash2,hash3"
	// @return 查询得到交易状态，比如 "status1,status2,status3"
	BatchFetchTransactionStatus(hashListString string) string

	// EstimateTransactionFee Most chains can estimate the fee directly to the transaction object
	// **But two chains don't work: `aptos`, `starcoin`**
	EstimateTransactionFee(transaction Transaction) (fee *OptionalString, err error)

	// EstimateTransactionFeeUsePublicKey All chains can call this method to estimate the gas fee.
	// **Chain  `aptos`, `starcoin` must pass in publickey**
	EstimateTransactionFeeUsePublicKey(transaction Transaction, pubkey string) (fee *OptionalString, err error)
}
