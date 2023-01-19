package utilHttp

import (
	"fmt"
	"testing"
)

func TestHttp(t *testing.T) {
	body := `{
	"token": "e492ce5094413d6122af53a22196dce3"
}`
	request, err := NewHttpRequest("https://xxxxx/openapi/xiaomi/news", "POST", body, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := DoHttpRequest(request, nil)
	if err != nil {
		fmt.Println(err)
	}
	respBody, _ := ReadRespBody(resp)
	fmt.Println(string(respBody))
}
