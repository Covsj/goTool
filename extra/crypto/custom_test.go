package crypto

import (
	"fmt"
	"testing"
)

func TestProcessKey(t *testing.T) {
	p := &CaesarUtil{
		OffsetFunc: func(str string) int {
			return len(str) % 3
		},
		StableIndexFunc: func(str string) map[int]bool {
			res := map[int]bool{}
			for i, _ := range str {
				if i != 0 && i != 5 && i%3 == 0 {
					res[i] = true
				}
			}
			return res
		},
		DupCountFunc: func(str string) int {
			return len(str) % 6
		},
	}
	ori := "b0W786t56iWc==6iy5BhdbWi68y6bw54guFghe6wC~"
	fmt.Println(ori)
	enc := p.DupEncodeKey(ori, true)
	fmt.Println(enc)
	dec := p.DupEncodeKey(enc, false)
	fmt.Println(dec)
	fmt.Println(ori == dec)
}
