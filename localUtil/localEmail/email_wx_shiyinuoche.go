package localEmail

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Covsj/goTool/localUtil"
)

var client = &http.Client{}

var TokenShiYiNuoChe = "juDjJIl8NMnwLyvwaEKR3zuy1N7kIyQs"

func SetNewEmail(accountLength int) string {
	//rand.Seed(time.Now().UnixNano())
	//pre := localUtil.GetMd5(time.Now().String())[:rand.Intn(3)+accountLength]
	//email := pre + "@mail.shiyinuoche.com"

	req, _ := http.NewRequest("POST", "https://mail.shiyinuoche.com/api/mailbox/rand_emprefix", nil)
	req.Header.Set("token", TokenShiYiNuoChe)
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
	res := RespWx{}
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		log.Error("设置邮箱失败 ", string(respBody))
		return ""
	}
	log.Info("设置邮箱返回 ", res.Data.Fulldomain)
	return res.Data.Fulldomain
}

func GetNewEmailEid(email string) string {
	request, _ := http.NewRequest("POST", "https://mail.shiyinuoche.com/api/mailbox/getnewest5", nil)
	request.Header.Add("token", TokenShiYiNuoChe)
	resp, err := client.Do(request)
	if err != nil {
		log.Error("获取验证码失败 Do ", err.Error())
		return ""
	}
	defer resp.Body.Close()
	res := WxShiYinNuoCheEmail{}
	body, err := localUtil.ReadRespBody(resp)
	if err != nil {
		log.Error("获取验证码失败  ReadRespBody ", err.Error())
		return ""
	}
	_ = json.Unmarshal(body, &res)
	log.Info("获取邮件长度  ", len(res.Data))

	for _, datum := range res.Data {
		if datum.To == email && strings.Contains(datum.Subject, "verification") {
			code := datum.Subject[len(datum.Subject)-6:]
			_, err := strconv.Atoi(code)
			if err == nil {
				return code
			}
		}
	}
	return ""
}
