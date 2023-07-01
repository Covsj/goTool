package localUtil

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Domain = "iubridge.com"
var DefaultId = "8e471dba9d6a1066f5a01b8e7c13b2cc"

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
		Date      int64       `json:"date"`
		Html      interface{} `json:"html"`
		MessageId string      `json:"messageId"`
		Name      string      `json:"name"`
		Eid       string      `json:"eid"`
	} `json:"data"`
	Status string `json:"status"`
}

func GetAEmailAccount() string {
	//DefaultId = GetMd5(time.Now().String())

	rand.Seed(time.Now().UnixNano())
	pre := GetMd5(time.Now().String())[:rand.Intn(3)+6]
	email := pre + "@" + Domain
	return email
}

// GetEid 传入邮箱名和想获取邮件发件人 返回邮箱最新邮件id
func GetEid(name string, senderName string) (Eid string, err error) {
	if strings.Contains(name, "@") {
		index := strings.Index(name, "@")
		name = name[:index]
	}
	url := "https://www.linshi-email.com/api/v1/refreshmessage/" + DefaultId + "/" + name + "@" + Domain + "?" +
		"t=" + strconv.Itoa(int(time.Now().Unix())) + "000"
	_, body, err := CallHttp(url, "GET", "", nil)
	if err != nil {
		return "", err
	}
	res := Resp{}
	_ = json.Unmarshal(body, &res)

	for i, list := range res.List {
		if list.AddressFrom == senderName {
			Eid = res.List[i].Eid
			return Eid, nil
		}
	}
	return "", err
}

// GetEmail 传入邮件Eid 获取邮件内容
func GetEmail(Eid string) (*Email, error) {
	url := "https://www.linshi-email.com/api/v1/getemailContent/" + Eid
	_, body, err := CallHttp(url, "GET", "", nil)
	if err != nil {
		return nil, err
	}
	res := &Email{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
