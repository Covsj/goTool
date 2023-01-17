package crypto

import (
	"fmt"
	"strings"
)

// Caesar 凯撒加密升级版
func Caesar(str string, isEncrypt bool) string {
	keyword := "HOMEPAGEDETAILEXPLORERELATEDuidnews" //加密密钥
	start := 6
	if isEncrypt {
		return VirginiaEncode(str, keyword, start)
	}
	return VirginiaDecode(str, keyword, start)
}

var numberMap = map[int32]int32{
	'1': '9',
	'9': '1',
	'4': '7',
	'7': '4',
	'6': '0',
	'0': '6',
	'3': '2',
	'2': '3',
}

func VirginiaEncode(msg, keyword string, start int) (message string) {
	msg = strings.ReplaceAll(msg, "&", "<")
	msg = strings.ReplaceAll(msg, "=", ">")
	msg = strings.ReplaceAll(msg, "-", "(")
	msg = strings.ReplaceAll(msg, "_", ")")

	keyNum := key2Value(keyword, start)

	for i, c := range msg {
		if c >= 'A' && c <= 'Z' {
			c += keyNum[i%len(keyNum)] // 轮换密钥应对值
			if c > 'Z' {
				c -= 26 // c 计算后不为大写字母则 -26
			}
		}
		if c >= 'a' && c <= 'z' {
			c += keyNum[i%len(keyNum)]
			if c > 'z' {
				c -= 26
			}
		}
		if value, ok := numberMap[c]; ok {
			c = value
		}
		message += fmt.Sprintf("%s", string(c))
	}

	return
}

func VirginiaDecode(msg, keyword string, start int) (message string) {
	msg = strings.ReplaceAll(msg, "<", "&")
	msg = strings.ReplaceAll(msg, ">", "=")
	msg = strings.ReplaceAll(msg, "(", "-")
	msg = strings.ReplaceAll(msg, ")", "_")

	keyNum := key2Value(keyword, start)
	length := len(keyNum)
	for i, c := range msg {
		if c >= 'A' && c <= 'Z' {
			c -= keyNum[i%length]
			if c < 'A' {
				c += 26
			}
		}
		if c >= 'a' && c <= 'z' {
			c -= keyNum[i%length]
			if c < 'a' {
				c += 26
			}
		}
		if value, ok := numberMap[c]; ok {
			c = value
		}
		message += fmt.Sprintf("%s", string(c))
	}

	return
}

func key2Value(keyword string, start int) []int32 {
	length := len(keyword)
	// 声明对应固定长度切片
	keyNum := make([]int32, length)
	for i := 0; i < length; i++ {
		for j := 0; j < 26; j++ {
			if uint8(j+65) == keyword[i] {
				keyNum[i] = int32(j + start)
				continue
			}
			if uint8(j+97) == keyword[i] {
				keyNum[i] = int32(j + start)
				continue
			}
		}
	}
	return keyNum
}
