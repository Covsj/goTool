package email

import (
	"fmt"
	"testing"
	"time"
)

func TestEmail(t *testing.T) {
	for {
		fmt.Println("正在获取验证码")
		Eid, err := GetEid("testEmail", "info@tokenview.io")
		if err != nil {
			fmt.Println(err)
		}
		if Eid != "" {
			email, err := GetEmail(Eid)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(email)
			break
		}
		time.Sleep(3 * time.Second)
	}
}
