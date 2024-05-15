package crypto

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	md5 := Md5("11")
	fmt.Println(md5)

	encode := Base64Encode(md5)
	fmt.Println(encode)
	res, err := Base64decode(encode)
	fmt.Println(res, err)
	fmt.Println(res == md5)
}
