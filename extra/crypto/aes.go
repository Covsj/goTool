package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

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
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
