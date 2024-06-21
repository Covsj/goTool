package basic

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/tyler-smith/go-bip39"
)

// MapListConcurrent 多协程转换数据
func MapListConcurrent(list []interface{}, threads int, fn func(interface{}) (interface{}, error)) ([]interface{}, error) {
	thread := 0
	max := threads
	wg := sync.WaitGroup{}

	mapContainer := newSafeMap()
	var firstError error
	for _, item := range list {
		if firstError != nil {
			continue
		}
		if max == 0 {
			wg.Add(1)
			// no limit
		} else {
			if thread == max {
				wg.Wait()
				thread = 0
			}
			if thread < max {
				wg.Add(1)
			}
		}

		go func(w *sync.WaitGroup, item interface{}, mapContainer *safeMap, firstError *error) {
			res, err := fn(item)
			if *firstError == nil && err != nil {
				*firstError = err
			} else {
				mapContainer.writeMap(item, res)
			}
			wg.Done()
		}(&wg, item, mapContainer, &firstError)
		thread++
	}
	wg.Wait()
	if firstError != nil {
		return nil, firstError
	}

	var result []interface{}
	for _, item := range list {
		result = append(result, mapContainer.Map[item])
	}
	return result, nil
}

// MapListConcurrentStringToString 多协程转换字符串数据,默认10协程
func MapListConcurrentStringToString(strList []string, fn func(string) (string, error)) ([]string, error) {
	list := make([]interface{}, len(strList))
	for i, s := range strList {
		list[i] = s
	}
	temp, err := MapListConcurrent(list, 10, func(i interface{}) (interface{}, error) {
		return fn(i.(string))
	})
	if err != nil {
		return nil, err
	}

	result := make([]string, len(temp))
	for i, v := range temp {
		result[i] = v.(string)
	}
	return result, nil
}

func MaxBigInt(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return x
	} else {
		return y
	}
}

func Max[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | string](x, y T) T {
	if x >= y {
		return x
	} else {
		return y
	}
}

func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | string](x, y T) T {
	if x <= y {
		return x
	} else {
		return y
	}
}

/* [zh] 该方法会捕捉 panic 抛出的值，并转成一个 error 对象通过参数指针返回
 *      注意: 如果想要返回它抓住的 error, 必须使用命名返回值！！
 * [en] This method will catch the value thrown by panic, and turn it into an error object and return it through the parameter pointer
 *		Note: If you want to return the error it caught, you must use a named return value! !
 *  ```
 *  func actionWillThrowError(parameters...) (namedErr error, other...) {
 *      defer CatchPanicAndMapToBasicError(&namedErr)
 *      // action code ...
 *      return namedErr, other...
 *  }
 *  ```
 */
func CatchPanicAndMapToBasicError(errOfResult *error) {
	// first we have to recover()
	errOfPanic := recover()
	if errOfResult == nil {
		return
	}
	if errOfPanic != nil {
		*errOfResult = MapAnyToBasicError(errOfPanic)
	} else {
		*errOfResult = MapAnyToBasicError(*errOfResult)
	}
}

func MapAnyToBasicError(e any) error {
	if e == nil {
		return nil
	}

	err, ok := e.(error)
	if ok {
		return errors.New(err.Error())
	}

	msg, ok := e.(string)
	if ok {
		return errors.New("panic error: " + msg)
	}

	code, ok := e.(int)
	if ok {
		return errors.New("panic error: code = " + strconv.Itoa(code))
	}

	return errors.New("panic error: unexpected error.")
}

// ParseNumber
// @param num any format number, such as decimal "1237890123", hex "0x123ef0", "123ef0"
func ParseNumber(num string) (*big.Int, error) {
	if strings.HasPrefix(num, "0x") || strings.HasPrefix(num, "0X") {
		num = num[2:]
		if b, ok := big.NewInt(0).SetString(num, 16); ok {
			return b, nil
		}
	}
	if b, ok := big.NewInt(0).SetString(num, 10); ok {
		return b, nil
	}
	if b, ok := big.NewInt(0).SetString(num, 16); ok {
		return b, nil
	}
	return nil, errors.New("invalid number")
}

// ParseNumberToHex
// @param num any format number, such as decimal "1237890123", hex "0x123ef0", "123ef0"
// @return hex number start with 0x, characters include 0-9 a-f
func ParseNumberToHex(num string) string {
	if b, err := ParseNumber(num); err == nil {
		return "0x" + b.Text(16)
	}
	return "0x0"
}

// ParseNumberToDecimal
// @param num any format number, such as decimal "1237890123", hex "0x123ef0", "123ef0"
// @return decimal number, characters include 0-9
func ParseNumberToDecimal(num string) string {
	if b, err := ParseNumber(num); err == nil {
		return b.Text(10)
	}
	return "0"
}

func BigIntMultiply(b *big.Int, ratio float64) *big.Int {
	f1 := new(big.Float).SetInt(b)
	product := f1.Mul(f1, big.NewFloat(ratio))
	res, _ := product.Int(big.NewInt(0))
	return res
}

// CalculateLastWord 根据传入的11个助记词，计算最后助记词
func CalculateLastWord(mnemonicWords []string) (string, error) {
	if len(mnemonicWords) != 11 {
		return "", errors.New("mnemonicWords not 11 length")
	}
	// found own morning
	wordList := bip39.GetWordList()
	m := []string{}
	for _, word := range wordList {
		mnemonic := fmt.Sprintf("%s %s", strings.Join(mnemonicWords, " "), word)
		if bip39.IsMnemonicValid(mnemonic) {
			m = append(m, word)
		}
	}

	if len(m) <= 0 {
		return "", errors.New("not found")
	}
	stringToNumber := func(str string, length int) int {
		h := fnv.New32a()
		_, err := h.Write([]byte(str))
		if err != nil {
			return 0
		}
		hashValue := h.Sum32()
		number := int(hashValue) % length
		return number
	}
	index := stringToNumber(mnemonicWords[len(mnemonicWords)-2]+mnemonicWords[len(mnemonicWords)-1], len(m))
	mnemonic := fmt.Sprintf("%s %s", strings.Join(mnemonicWords, " "), m[index])
	return mnemonic, nil
}
