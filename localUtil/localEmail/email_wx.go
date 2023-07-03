package localEmail

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Covsj/goTool/localUtil"
)

var client = &http.Client{}

type VxEmail struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Host   string `json:"host"`
	Domain string `json:"domain"`
}

func NewWxEmail(name, token, host, domain string) *VxEmail {
	return &VxEmail{
		Name:   name,
		Token:  token,
		Host:   host,
		Domain: domain,
	}
}

func NewEmailTools() []*VxEmail {
	res := []*VxEmail{
		&VxEmail{
			Name:   "邮箱快捷助手",
			Token:  "juDjJIl8NMnwLyvwaEKR3zuy1N7kIyQs",
			Host:   "mail.shiyinuoche.com",
			Domain: "mail.shiyinuoche.com",
		},
		&VxEmail{
			Name:   "随机邮箱",
			Token:  "OVYPbGekfMBmQkCIP42VMRfq4scAwG09",
			Host:   "mail.11811.cn",
			Domain: "0dg.top",
		},
		&VxEmail{
			Name:   "个性随机邮箱",
			Token:  "s27PlGt7op1NEvVdVm3PNTslLId1dEjR",
			Host:   "email.1718u.com",
			Domain: "email.nm.cn",
		},
		&VxEmail{
			Name:   "邮箱",
			Token:  "zSh1bZL8yCMLdyNgqIHGCgbvEMw4Uc4t",
			Host:   "mail.pianyueniao.top",
			Domain: "pianyueniao.top",
		},
	}
	return res
}

func (cfg *VxEmail) SetNewWxEmail(accountLength int) string {
	var body string
	if accountLength != 0 {
		rand.Seed(time.Now().UnixNano())
		pre := localUtil.GetMd5(time.Now().String())[:rand.Intn(3)+accountLength]

		marshal, err := json.Marshal(&wxEmailSetEmailRep{EmPrefix: pre})
		if err != nil {
			log.Error("设置邮箱失败 反序列化失败 ", err.Error())
		} else {
			body = string(marshal)
		}
	}

	req, _ := http.NewRequest("POST", "https://"+cfg.Domain+"/api/mailbox/rand_emprefix",
		strings.NewReader(body))
	req.Header.Set("token", cfg.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("设置邮箱失败 ", err.Error())
		return ""
	}
	respBody, err := localUtil.ReadRespBody(resp)
	if err != nil {
		log.Error("设置邮箱失败 读取body失败 ", err.Error())
		return ""
	}
	res := wxEmailSetResp{}
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		log.Error("设置邮箱失败 ", string(respBody))
		return ""
	}
	log.Info("设置邮箱返回 ", res.Data.Fulldomain)
	return res.Data.Fulldomain
}

func (cfg *VxEmail) GetEmailMsgFromWx(email string) []WxEmailRespData {
	var result []WxEmailRespData

	request, _ := http.NewRequest("POST", "https://"+cfg.Host+"/api/mailbox/getnewest5", nil)
	request.Header.Add("token", cfg.Token)
	resp, err := client.Do(request)
	if err != nil {
		log.Error("获取验证码失败 Do ", err.Error())
		return result
	}
	defer resp.Body.Close()
	res := wxEmailGetMsgResp{}
	body, err := localUtil.ReadRespBody(resp)
	if err != nil {
		log.Error("获取验证码失败  ReadRespBody ", err.Error())
		return result
	}
	_ = json.Unmarshal(body, &res)
	log.Info("获取邮件长度  ", len(res.Data))

	for _, datum := range res.Data {
		if datum.To == email {
			result = append(result, datum)
		}
	}
	return result
}
