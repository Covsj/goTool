package eth

import (
	"fmt"
	"testing"
)

func TestTransformEIP55Address(t *testing.T) {
	oriAddress := "0xf98d6a61951ee2d6d1ee5debc9b3758a88888888"
	eIP55Address := TransformEIP55Address(oriAddress)
	fmt.Println(oriAddress)
	fmt.Println(eIP55Address)
	fmt.Println(oriAddress == eIP55Address)
}
