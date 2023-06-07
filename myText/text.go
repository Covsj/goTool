package myText

import "strings"

func GetMediateText(str, before, after string) string {
	i := strings.Index(str, before)
	j := strings.Index(str, after)
	if i < 0 || j < 0 {
		return ""
	}
	return str[i+len(before) : j]
}
