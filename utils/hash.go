package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

var CaesarKeyWord = "ABedHswOMEPAGEaDETAIdeLEXsPLORdfeERELATEDuidnews"

func GetMd5(ori string) string {
	h := md5.New()
	h.Write([]byte(ori))
	toString := hex.EncodeToString(h.Sum(nil))
	return strings.ToLower(toString)
}

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

// Caesar 凯撒加密升级版
func Caesar(str string, isEncrypt bool) string {
	keyword := CaesarKeyWord //加密密钥
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

// AesEncryptCBC Aes加密 aesKey 16, 24或者32
func AesEncryptCBC(origStr string, aesKey []byte) (string, error) {
	origData := []byte(origStr)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()                                 // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                   // 补全码
	blockMode := cipher.NewCBCEncrypter(block, aesKey[:blockSize]) // 加密模式
	encrypted := make([]byte, len(origData))                       // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                     // 加密
	toString := base64.StdEncoding.EncodeToString(encrypted)
	return toString, nil
}

// AesDecryptCBC Aes解密
func AesDecryptCBC(origStr string, aesKey []byte) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(origStr)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()                                 // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, aesKey[:blockSize]) // 加密模式
	decrypted := make([]byte, len(encrypted))                      // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                    // 解密
	decrypted = pkcs5UnPadding(decrypted)                          // 去除补全码
	return string(decrypted), nil
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
