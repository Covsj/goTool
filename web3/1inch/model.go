package _inch

type Web3Call struct {
	Data     string `json:"data"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gasLimit"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

type TokenAllowance struct {
	Allowance string `json:"allowance"`
}

type SwapParam struct {
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
	Amount           string `json:"amount"`
	FromAddress      string `json:"fromAddress"`
	Slippage         int    `json:"slippage"`
	DisableEstimate  bool   `json:"disableEstimate"`
	AllowPartialFill bool   `json:"allowPartialFill"`
}
