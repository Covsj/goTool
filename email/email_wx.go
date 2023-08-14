package email

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Covsj/goTool/utils"
	log "github.com/sirupsen/logrus"
)

var htppClient = &http.Client{}

type WxEmail struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Host   string `json:"host"`
	Domain string `json:"domain"`
}

var List []*WxEmail

func init() {
	List = []*WxEmail{
		&WxEmail{
			Name:   "邮箱快捷助手",
			Token:  "juDjJIl8NMnwLyvwaEKR3zuy1N7kIyQs",
			Host:   "mail.shiyinuoche.com",
			Domain: "mail.shiyinuoche.com",
		},
		&WxEmail{
			Name:   "随机邮箱",
			Token:  "OVYPbGekfMBmQkCIP42VMRfq4scAwG09",
			Host:   "mail.11811.cn",
			Domain: "0dg.top",
		},
		&WxEmail{
			Name:   "个性随机邮箱",
			Token:  "s27PlGt7op1NEvVdVm3PNTslLId1dEjR",
			Host:   "email.1718u.com",
			Domain: "email.nm.cn",
		},
		&WxEmail{
			Name:   "邮箱",
			Token:  "zSh1bZL8yCMLdyNgqIHGCgbvEMw4Uc4t",
			Host:   "mail.pianyueniao.top",
			Domain: "pianyueniao.top",
		},
		&WxEmail{
			Name:   "快趣邮箱",
			Token:  "fU1egC20e93jZgvvwkxZASEWPBn6bWWQ",
			Host:   "yx.kuaiquyun.cn",
			Domain: "yx.kuaiquyun.cn",
		},
		&WxEmail{
			Name:   "极速邮箱管家",
			Token:  "TRe8Q6jQCVh0O6HS2r9tmRna2vqkIDSe",
			Host:   "x.chuanyueshikong.cloud",
			Domain: "shenmak.cn",
		},
		&WxEmail{
			Name:   "无限邮箱",
			Token:  "BHlsuXDYzT2L9l6G083nXEoRXTkaCJAI",
			Host:   "mail.ynnmx.com",
			Domain: "lemonclo.com",
		},
		&WxEmail{
			Name:   "快邮箱",
			Token:  "GMxrUqbL9APnCBLHdlRqWamMBCc3Cavu",
			Host:   "mail.fyxkpro.cn",
			Domain: "mail.fyxkpro.cn",
		},
		&WxEmail{
			Name:   "微邮箱",
			Token:  "AuRnqEbxcxMICtoRDhE6La3VhICffAb5",
			Host:   "yx.mageline.cc",
			Domain: "wyx.xjcms.cc",
		},
		&WxEmail{
			Name:   "万能邮箱",
			Token:  "lY54N67b5JRy9d4rYurp3n8pqYEcVIFS",
			Host:   "yx.nuotui.vip",
			Domain: "yx.nuotui.vip",
		},
		&WxEmail{
			Name:   "临时邮箱生成器",
			Token:  "zuUyZm6utiKnHvQw0WHtOkSHzFUFCYGg",
			Host:   "mail.qizs.cn",
			Domain: "colingpt.pro",
		},
		&WxEmail{
			Name:   "个人无限邮箱",
			Token:  "nEGeek0P1Bg4MKP8vgjdJJregcjZXNPC",
			Host:   "mail.picyw.com",
			Domain: "52mail.cab",
		},
		&WxEmail{
			Name:   "我的个人邮箱",
			Token:  "4mB1xvGMLO80CrmFT42oU40FVjF0NjPh",
			Host:   "mail.sihong.vip",
			Domain: "156.email",
		},
	}
}

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
	resp, err := htppClient.Do(req)
	if err != nil {
		log.Error("设置邮箱失败 ", err.Error())
		return ""
	}
	respBody, err := utils.ReadRespBody(resp)
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
	resp, err := htppClient.Do(request)
	if err != nil {
		log.Error("获取验证码失败 Do ", err.Error())
		return result
	}
	defer resp.Body.Close()
	res := wxEmailGetMsgResp{}
	body, err := utils.ReadRespBody(resp)
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
