package crypto

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	md5 := HashGenerator("hello world", "Sha512-224")
	fmt.Println(md5.ToHexString())

	encode := BaseEncode(md5.ToHexString(), "100")
	fmt.Println(encode.ToString())
	decode := BaseDecode(encode.String(), "100")
	fmt.Println(decode.ToString())
	fmt.Println(decode.ToString() == md5.ToHexString())
}
