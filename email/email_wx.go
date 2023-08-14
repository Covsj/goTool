package email

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Covsj/goTool/httpTool"
	"github.com/Covsj/goTool/utils"
	log "github.com/sirupsen/logrus"
)

var httpClient = &http.Client{}

type WxEmail struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Host   string `json:"host"`
	Domain string `json:"domain"`
}

var List []*WxEmail

func NewWxEmail(name, token, host, domain string) *WxEmail {
	if name != "" && token != "" && host != "" && domain != "" {
		return &WxEmail{
			Name:   name,
			Token:  token,
			Host:   host,
			Domain: domain,
		}
	}
	rand.Seed(time.Now().UnixNano())
	return List[rand.Intn(len(List))]
}

func (cfg *WxEmail) SetNewWxEmail(accountLength int) string {
	var body string
	if accountLength != 0 {
		rand.Seed(time.Now().UnixNano())
		pre := utils.GetMd5(time.Now().String())[:rand.Intn(3)+accountLength]

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
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error("设置邮箱失败 ", err.Error())
		return ""
	}
	respBody, err := httpTool.ReadRespBody(resp)
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

func (cfg *WxEmail) GetEmailMsgFromWx(email string) []WxEmailRespData {
	var result []WxEmailRespData

	request, _ := http.NewRequest("POST", "https://"+cfg.Host+"/api/mailbox/getnewest5", nil)
	request.Header.Add("token", cfg.Token)
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Error("获取验证码失败 Do ", err.Error())
		return result
	}
	defer resp.Body.Close()
	res := wxEmailGetMsgResp{}
	body, err := httpTool.ReadRespBody(resp)
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
