package email

import (
	"encoding/json"
	"strconv"
	"time"

	"goTool/utilHttp"
)

var domain = "iubridge.com"

type Resp struct {
	List []struct {
		NameTo      string `json:"name_to"`
		NameFrom    string `json:"name_from"`
		Eid         string `json:"eid"`
		ESubject    string `json:"e_subject"`
		EDate       int64  `json:"e_date"`
		AddressFrom string `json:"address_from"`
	} `json:"list"`
	Status string `json:"status"`
}

type Email struct {
	Data struct {
		To      string `json:"to"`
		Seqno   int    `json:"seqno"`
		Subject string `json:"subject"`
		From    struct {
			Address string `json:"address"`
			Name    string `json:"name"`
		} `json:"from"`
		Date      int64  `json:"date"`
		Html      string `json:"html"`
		MessageId string `json:"messageId"`
		Name      string `json:"name"`
		Eid       string `json:"eid"`
	} `json:"data"`
	Status string `json:"status"`
}

// GetEid 传入邮箱名和想获取邮件发件人 返回邮箱最新邮件id
func GetEid(name string, senderName string) (Eid string, err error) {
	url := "https://www.linshi-email.com/api/v1/refreshmessage/61b6348b42316c1967348fbc347b3c70/" + name + "@" + domain + "?" +
		"t=" + strconv.Itoa(int(time.Now().Unix())) + "000"
	resp, err := utilHttp.CallHttp(url, "GET", "", nil)
	if err != nil {
		return "", err
	}
	res := Resp{}
	_ = json.Unmarshal(resp, &res)

	for i, list := range res.List {
		if list.AddressFrom == senderName {
			Eid = res.List[i].Eid
			return Eid, nil
		}
	}
	return "", err
}

// GetEmail 传入邮件Eid 获取邮件内容
func GetEmail(Eid string) (string, error) {
	url := "https://www.linshi-email.com/api/v1/getemailContent/" + Eid
	resp, err := utilHttp.CallHttp(url, "GET", "", nil)
	if err != nil {
		return "", err
	}
	res := Email{}
	_ = json.Unmarshal(resp, &res)
	e := res.Data.Html
	return e, nil
}
