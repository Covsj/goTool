package _inch

import (
	"encoding/json"

	gotool_http "github.com/Covsj/goTool/http"
)

func Get1inchCheckAllowance(chainId, tokenAddress, walletAddress string, url string) (string, error) {
	if url == "" {
		url = getBaseApiUrl(chainId)
	}
	url += "/approve/allowance?tokenAddress=" + tokenAddress + "&walletAddress=" + walletAddress
	_, body, err := gotool_http.Send(&gotool_http.RequestOptions{URL: url})
	//fmt.Println(resp, string(body))
	if err != nil {
		return "", err
	}
	var m TokenAllowance
	err = json.Unmarshal(body, &m)
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
	_, body, err := gotool_http.Send(&gotool_http.RequestOptions{URL: url})
	if err != nil {
		return nil, err
	}
	m := &Web3Call{GasLimit: "200000"}
	err = json.Unmarshal(body, m)
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
