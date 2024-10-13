package _inch

import (
	gotool_http "github.com/Covsj/goTool/http"
)

func Get1inchCheckAllowance(chainId, tokenAddress, walletAddress string, url string) (string, error) {
	if url == "" {
		url = getBaseApiUrl(chainId)
	}
	m := &TokenAllowance{}
	url += "/approve/allowance?tokenAddress=" + tokenAddress + "&walletAddress=" + walletAddress
	_, err := gotool_http.DoRequest(&gotool_http.ReqOpt{Url: url, RespOut: m})
	//fmt.Println(resp, string(body))
	if err != nil {
		return "", err
	}
	return m.Allowance, nil
}

func Approve1inchToken(chainId, tokenAddress, amount string) (*Web3Call, error) {
	url := getBaseApiUrl(chainId)
	url += "/approve/transaction?tokenAddress=" + tokenAddress
	if amount != "" {
		url += "&amount=" + amount
	}
	m := &Web3Call{GasLimit: "200000"}
	_, err := gotool_http.DoRequest(&gotool_http.ReqOpt{Url: url, RespOut: m})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func getBaseApiUrl(chainId string) string {
	return "https://api.1inch.io/v5.0/" + chainId
}
func getBroadcastApiUrl(chainId string) string {
	return "https://tx-gateway.1inch.io/v1.1/" + chainId + "/broadcast"
}
