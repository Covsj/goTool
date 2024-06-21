package eth

const (
	Erc20MethodTransfer  = "transfer"
	Erc20MethodApprove   = "approve"
	Erc20MethodBalanceOf = "balanceOf"
	Erc20MethodDecimals  = "decimals"
)

// 默认gas limit估算失败后，21000 * 3 = 63000
const (
	DEFAULT_CONTRACT_GAS_LIMIT = "63000"
	DEFAULT_ETH_GAS_LIMIT      = "21000"
	// 当前网络 standard gas price
	DEFAULT_ETH_GAS_PRICE = "20000000000"
)
