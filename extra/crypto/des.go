package crypto

import (
	"gitee.com/golang-module/dongle"
	"strings"
)

func DesEncryptData(data interface{}, mode, padding string, desKey, desIv interface{}) dongle.Encrypter {
	cipher := NewCipher(mode, padding, desKey, desIv)
	var encrypter dongle.Encrypter
	switch v := data.(type) {
	case string:
		encrypter = dongle.Encrypt.FromString(v).ByDes(cipher)
	case []byte:
		encrypter = dongle.Encrypt.FromBytes(v).ByDes(cipher)
	default:
		return encrypter
	}
	return encrypter
}

func DesDecryptData(encryptedData interface{}, mode, padding string, desKey, desIv interface{}, encodingMode string) dongle.Decrypter {

	cipher := NewCipher(mode, padding, desKey, desIv)

	encodingMode = strings.ToLower(encodingMode)
	switch encodingMode {
	case "raw":
		return dongle.Decrypt.FromRawString(encryptedData.(string)).ByDes(cipher)
	case "hex":
		return dongle.Decrypt.FromHexString(encryptedData.(string)).ByDes(cipher)
	case "base64":
		return dongle.Decrypt.FromBase64String(encryptedData.(string)).ByDes(cipher)
	case "bytes":
		return dongle.Decrypt.FromRawBytes(encryptedData.([]byte)).ByDes(cipher)
	case "hex-bytes":
		return dongle.Decrypt.FromHexBytes(encryptedData.([]byte)).ByDes(cipher)
	case "base64-bytes":
		return dongle.Decrypt.FromBase64Bytes(encryptedData.([]byte)).ByDes(cipher)
	}

	return dongle.Decrypter{}
}

func Des3EncryptData(data interface{}, mode, padding string, aesKey, aesIv interface{}) dongle.Encrypter {
	cipher := NewCipher(mode, padding, aesKey, aesIv)
	var encrypter dongle.Encrypter
	switch v := data.(type) {
	case string:
		encrypter = dongle.Encrypt.FromString(v).By3Des(cipher)
	case []byte:
		encrypter = dongle.Encrypt.FromBytes(v).By3Des(cipher)
	default:
		return encrypter
	}
	return encrypter
}

func Des3DecryptData(encryptedData interface{}, mode, padding string, aesKey, aesIv interface{}, encodingMode string) dongle.Decrypter {

	cipher := NewCipher(mode, padding, aesKey, aesIv)

	encodingMode = strings.ToLower(encodingMode)
	switch encodingMode {
	case "bytes":
		return dongle.Decrypt.FromRawBytes(encryptedData.([]byte)).By3Des(cipher)
	case "hex-bytes":
		return dongle.Decrypt.FromHexBytes(encryptedData.([]byte)).By3Des(cipher)
	case "base64-bytes":
		return dongle.Decrypt.FromBase64Bytes(encryptedData.([]byte)).By3Des(cipher)
	case "hex":
		return dongle.Decrypt.FromHexString(encryptedData.(string)).By3Des(cipher)
	case "base64":
		return dongle.Decrypt.FromBase64String(encryptedData.(string)).By3Des(cipher)
	default:
		return dongle.Decrypt.FromRawString(encryptedData.(string)).By3Des(cipher)
	}
}
