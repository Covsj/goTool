package crypto

import (
	"crypto/md5"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func Md5(ori string) string {
	h := md5.New()
	h.Write([]byte(ori))
	toString := hex.EncodeToString(h.Sum(nil))
	return strings.ToLower(toString)
}

func Base64Encode(ori string) string {
	input := []byte(ori)
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func Base64decode(ori string) (res string, err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(ori)
	if err != nil {
		return "", err
	}
	return string(decodeBytes), nil
}

func Base32Encode(ori string) string {
	input := []byte(ori)
	encodeString := base32.StdEncoding.EncodeToString(input)
	return encodeString
}

func Base32decode(ori string) (res string, err error) {
	decodeBytes, err := base32.StdEncoding.DecodeString(ori)
	if err != nil {
		return "", err
	}
	return string(decodeBytes), nil
}
