package localEmail

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Covsj/goTool/localUtil"
)

var DefaultIdLinShi = "61b6348b42316c1967348fbc347b3c70"

func GetAEmailAccount(accountLength int) string {
	if accountLength == 0 {
		accountLength = 6
	}
	rand.Seed(time.Now().UnixNano())
	pre := localUtil.GetMd5(time.Now().String())[:rand.Intn(3)+accountLength]
	email := pre + "@iubridge.com"
	return email
}

// GetEid 传入邮箱名和想获取邮件发件人 返回邮箱最新邮件id
func GetEid(name string, senderName string) (Eid string, err error) {
	if strings.Contains(name, "@") {
		index := strings.Index(name, "@")
		name = name[:index]
	}
	url := "https://www.linshi-email.com/api/v1/refreshmessage/" + DefaultIdLinShi + "/" + name + "@iubridge.com?" +
		"t=" + strconv.Itoa(int(time.Now().Unix())) + "000"
	_, body, err := localUtil.CallHttp(url, "GET", "", nil)
	if err != nil {
		return "", err
	}
	res := LinShiEmailResp{}
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
func GetEmail(Eid string) (*LinShiEmail, error) {
	url := "https://www.linshi-email.com/api/v1/getemailContent/" + Eid
	_, body, err := localUtil.CallHttp(url, "GET", "", nil)
	if err != nil {
		return nil, err
	}
	res := &LinShiEmail{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetEmailAll(name string, senderName string) *LinShiEmail {
	for {
		time.Sleep(1 * time.Second)
		eid, err := GetEid(name, senderName)
		if err != nil {
			continue
		}
		if eid == "" {
			continue
		}
		email, err := GetEmail(eid)
		if err != nil {
			continue
		}
		return email
	}
}
