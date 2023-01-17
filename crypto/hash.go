package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func GetMd5(ori string) string {
	h := md5.New()
	h.Write([]byte(ori))
	toString := hex.EncodeToString(h.Sum(nil))
	return strings.ToLower(toString)
}
