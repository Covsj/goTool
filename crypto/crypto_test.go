package crypto

import (
	"fmt"
	"testing"
)

func TestCaesar(t *testing.T) {
	oriStr := "yongganniuniu"
	crypto := Caesar(oriStr, true)
	fmt.Println(crypto)
	fmt.Println(Caesar(crypto, false))
}
