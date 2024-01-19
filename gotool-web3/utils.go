package gotool_web3

import "fmt"

// 查找单词在列表中的索引
func findIndex(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// 二进制字符串转换为十进制
func binaryToDecimal(binaryStr string) (int, error) {
	var num int
	_, err := fmt.Sscanf(binaryStr, "%b", &num)
	return num, err
}
