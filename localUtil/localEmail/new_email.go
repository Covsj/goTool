package localEmail

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Aof/model"

	log "github.com/sirupsen/logrus"

	"github.com/Covsj/goTool/localUtil"
)

var token = ""
var list = []string{
	"uuf.me",
	"qabq.com",
	"nqmo.com",
	"end.tw",
	"uuf.me",
	"yzm.de",
}

func SetNewEmail(accountLength int) string {
	if accountLength == 0 {
		accountLength = 6
	}
	rand.Seed(time.Now().UnixNano())
	pre := localUtil.GetMd5(time.Now().String())[:rand.Intn(3)+accountLength]
	email := pre + "@" + list[rand.Intn(len(list))]

	request, _ := http.NewRequest("POST", "https://api.mail.cx/api/v1/auth/authorize_token", nil)

	resp, err := proCl.Do(request)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := localUtil.ReadRespBody(resp)
	if err != nil {
		return ""
	}

	token = strings.ReplaceAll(string(body), "\"", "")
	token = strings.ReplaceAll(token, "\n", "")
	log.Info(token)
	return email
}

func GetNewEmailEid(email string) string {
	reqUrl := "https://mail.cx/api/api/v1/mailbox/" + email
	request, _ := http.NewRequest("GET", reqUrl, nil)
	request.Header.Add("authorization", "bearer "+token)
	//log.Info(email, "   ", token)
	resp, err := proCl.Do(request)
	if err != nil {
		log.Error("获取验证码失败 Do ", err.Error())
		return ""
	}
	defer resp.Body.Close()
	res := []*model.NewEmailResp{}
	body, err := localUtil.ReadRespBody(resp)
	if err != nil {
		log.Error("获取验证码失败  ReadRespBody ", err.Error())
		return ""
	}
	_ = json.Unmarshal(body, &res)
	for _, datum := range res {
		if datum.From == "AOF Games <aof@aof.games>" {
			if strings.Contains(datum.Subject, "verification code") {
				code := datum.Subject[len(datum.Subject)-6:]
				_, err := strconv.Atoi(code)
				if err == nil {
					return code
				}
			}
		}
	}
	return ""
}
