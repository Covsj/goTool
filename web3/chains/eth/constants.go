package eth

const (
	Erc20MethodTransfer  = "transfer"
	Erc20MethodApprove   = "approve"
	Erc20MethodBalanceOf = "balanceOf"
	Erc20MethodDecimals  = "decimals"
)

// 默认gas limit估算失败后，21000 * 3 = 63000
const (
	DefaultContractGasLimit = "63000"
	DefaultEthGasLimit      = "21000"
	// 当前网络 standard gas price
	DefaultEthGasPrice = "20000000000"
)
