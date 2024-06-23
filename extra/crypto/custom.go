package crypto

import (
	"os"
	"strings"
)

type CaesarUtil struct {
	OffsetFunc      func(str string) int
	StableIndexFunc func(str string) map[int]bool
	CustomCharset   string
	DupCountFunc    func(str string) int
}

var (
	DefaultCustomCharset = os.Getenv("DefaultCustomCharset")
	DefaultOffset        = 13
)

func (p *CaesarUtil) EncodeKey(str string, encrypt bool) string {
	if !encrypt {
		str = reverseString(str)
	}

	offset := DefaultOffset
	if p.OffsetFunc != nil {
		offset = p.OffsetFunc(str)
	}
	if offset == 0 {
		offset = DefaultOffset
	}

	customCharset := p.CustomCharset
	if customCharset == "" {
		customCharset = DefaultCustomCharset
	}

	aes, err := AesEncryptCBC(Base64Encode(Base32Encode(customCharset)), []byte(Md5(customCharset)))
	if err != nil {
		aes = customCharset
	}
	customCharset = uniqueChars(customCharset + aes)

	stableIndex := map[int]bool{}
	if p.StableIndexFunc != nil {
		stableIndex = p.StableIndexFunc(str)
	}

	mapping := createMapping(customCharset, offset, encrypt)

	var result strings.Builder
	for i, char := range str {
		if stableIndex[i] {
			result.WriteRune(char)
			continue
		}
		if mappedChar, ok := mapping[string(char)]; ok {
			result.WriteString(mappedChar)
		} else {
			result.WriteRune(char)
		}
	}

	if encrypt {
		return reverseString(result.String())
	}

	return result.String()
}

func (p *CaesarUtil) DupEncodeKey(str string, encrypt bool) string {
	result := str
	var dupCount int
	if p.DupCountFunc != nil {
		dupCount = p.DupCountFunc(str)
	}
	if dupCount == 0 {
		dupCount = 1
	}
	for i := 0; i < dupCount; i++ {
		result = p.EncodeKey(result, encrypt)
	}
	return result
}

func reverseString(str string) string {
	runes := []rune(str)
	length := len(runes)

	switch {
	case length <= 10:
		return reverse(runes, 0, length-1)
	case length <= 20:
		return segmentedReverse(runes, 3)
	case length <= 30:
		return segmentedReverse(runes, 4)
	default:
		return segmentedReverse(runes, 5)
	}
}

func reverse(runes []rune, start, end int) string {
	for i, j := start, end; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes[start : end+1])
}

func segmentedReverse(runes []rune, segments int) string {
	length := len(runes)
	chunkSize := length / segments
	var result strings.Builder
	for i := 0; i < segments; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		if i == segments-1 {
			end = length - 1
		}
		result.WriteString(reverse(runes, start, end))
	}
	return result.String()
}

func uniqueChars(str string) string {
	seen := make(map[rune]bool)
	var result strings.Builder
	for _, char := range str {
		if !seen[char] {
			seen[char] = true
			result.WriteRune(char)
		}
	}
	return result.String()
}

func createMapping(charset string, offset int, encrypt bool) map[string]string {
	mapping := make(map[string]string)
	length := len(charset)
	for i, r := range charset {
		if encrypt {
			mapping[string(r)] = string(charset[(i+offset)%length])
		} else {
			mapping[string(charset[(i+offset)%length])] = string(r)
		}
	}
	return mapping
}
