package crypto

import (
	"encoding/base64"
	"strconv"
)

func GetBase64Encode(ori string) string {
	input := []byte(ori)
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func GetBase64decode(ori string) (res string, err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(ori)
	if err != nil {
		return "", err
	}
	return string(decodeBytes), nil
}

// Enc XorKey 8位
func Enc(src string, XorKey []byte) string {
	var result string
	j := 0
	s := ""
	bt := []rune(src)
	for i := 0; i < len(bt); i++ {
		s = strconv.FormatInt(int64(byte(bt[i])^XorKey[j]), 16)
		if len(s) == 1 {
			s = "0" + s
		}
		result = result + (s)
		j = (j + 1) % 8
	}
	return result
}

// Dec XorKey 8位
func Dec(src string, XorKey []byte) string {
	var result string
	var s int64
	j := 0
	bt := []rune(src)
	for i := 0; i < len(src)/2; i++ {
		s, _ = strconv.ParseInt(string(bt[i*2:i*2+2]), 16, 0)
		result = result + string(byte(s)^XorKey[j])
		j = (j + 1) % 8
	}
	return result
}
