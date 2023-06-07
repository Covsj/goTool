package dex

import (
	"encoding/json"

	"goTool/model"
	"goTool/myHttp"
)

func getBaseApiUrl(chainId string) string {
	return "https://api.1inch.io/v5.0/" + chainId
}
func getBroadcastApiUrl(chainId string) string {
	return "https://tx-gateway.1inch.io/v1.1/" + chainId + "/broadcast"
}

func CheckAllowance(chainId, tokenAddress, walletAddress string, url string) (string, error) {
	if url == "" {
		url = getBaseApiUrl(chainId)
	}
	url += "/approve/allowance?tokenAddress=" + tokenAddress + "&walletAddress=" + walletAddress
	resp, err := myHttp.CallHttp(url, "GET", "", nil)
	if err != nil {
		return "", err
	}

	var m model.CheckAllowance
	err = json.Unmarshal(resp, &m)
	if err != nil {
		return "", err
	}
	return m.Allowance, nil
}

func GetApproveTransaction(chainId, tokenAddress, amount string) (*model.Web3Call, error) {
	url := getBaseApiUrl(chainId)
	url += "/approve/transaction?tokenAddress=" + tokenAddress
	if amount != "" {
		url += "&amount=" + amount
	}
	resp, err := myHttp.CallHttp(url, "GET", "", nil)
	if err != nil {
		return nil, err
	}
	m := &model.Web3Call{GasLimit: "200000"}
	err = json.Unmarshal(resp, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
