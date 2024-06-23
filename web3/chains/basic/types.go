package basic

import (
	"sync"

	"github.com/Covsj/goTool/web3/chains/basic/inter"
)

type SDKEnumInt = int
type SDKEnumString = string

type OptionalString struct {
	Value string
}

func NewOptionalString(s string) *OptionalString {
	return &OptionalString{Value: s}
}

type OptionalBool struct {
	Value bool
}

func NewOptionalBool(b bool) *OptionalBool {
	return &OptionalBool{Value: b}
}

type OptionalInt struct {
	Value int
}

func NewOptionalInt(i int) *OptionalInt {
	return &OptionalInt{Value: i}
}

type safeMap struct {
	sync.RWMutex
	Map map[interface{}]interface{}
}

func newSafeMap() *safeMap {
	return &safeMap{Map: make(map[interface{}]interface{})}
}

func (l *safeMap) readMap(key interface{}) (interface{}, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *safeMap) writeMap(key interface{}, value interface{}) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

type StringArray struct {
	inter.AnyArray[string]
}

func NewStringArray() *StringArray {
	return &StringArray{[]string{}}
}

func NewStringArrayWithItem(elem string) *StringArray {
	return &StringArray{[]string{elem}}
}

func (a StringArray) Contains(value string) bool {
	idx := inter.FirstIndexOf(a.AnyArray, func(elem string) bool { return elem == value })
	return idx != -1
}

type StringMap struct {
	inter.AnyMap[string, string]
}

func NewStringMap() *StringMap {
	return &StringMap{map[string]string{}}
}

func (m *StringMap) Keys() *StringArray {
	keys := inter.KeysOf(m.AnyMap)
	return &StringArray{keys}
}
