package cache

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"sync"
	"time"
)

type CacheShard struct {
	MaxEntries int

	shardNum int
	lrus     []*Cache
	mus      []sync.Mutex
}

type TTLValue struct {
	expireTs int64
	value    interface{}
}

func NewCacheShard(maxEntries int, shardNum int) *CacheShard {
	if shardNum <= 1 {
		return nil
	}
	cacheShard := &CacheShard{
		MaxEntries: maxEntries,
		shardNum:   shardNum,
	}
	cacheShard.mus = make([]sync.Mutex, shardNum)
	cacheShard.lrus = make([]*Cache, shardNum)
	for i := 0; i < shardNum; i++ {
		cacheShard.lrus[i] = New(maxEntries / shardNum)
	}
	return cacheShard
}

func (c *CacheShard) getIndex(key Key) int {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	h := fnv.New32a()
	err := enc.Encode(key)
	if err != nil {
		return 0
	}
	_, _ = h.Write(buf.Bytes())
	return (int)(h.Sum32()) % c.shardNum
}

func (c *CacheShard) Add(key Key, value interface{}) {
	index := c.getIndex(key)
	c.mus[index].Lock()
	defer c.mus[index].Unlock()
	ttlValue := &TTLValue{value: value}
	c.lrus[index].Add(key, ttlValue)
}

func (c *CacheShard) TtlAdd(key Key, value interface{}, ttl time.Duration) {
	index := c.getIndex(key)
	c.mus[index].Lock()
	defer c.mus[index].Unlock()
	ttlValue := &TTLValue{expireTs: time.Now().Add(ttl).Unix(), value: value}
	c.lrus[index].Add(key, ttlValue)

}

func (c *CacheShard) Get(key Key) (value interface{}, ok bool) {
	index := c.getIndex(key)
	c.mus[index].Lock()
	defer c.mus[index].Unlock()
	if v, ok := c.lrus[index].Get(key); ok {
		ttlValue := v.(*TTLValue)
		if ttlValue.expireTs > 0 && ttlValue.expireTs < time.Now().Unix() {
			c.lrus[index].Remove(key)
			return nil, false
		}
		return ttlValue.value, true
	}
	return nil, false
}

func (c *CacheShard) Remove(key Key) {
	index := c.getIndex(key)
	c.mus[index].Lock()
	defer c.mus[index].Unlock()
	c.lrus[index].Remove(key)
}

func (c *CacheShard) RemoveOldest() {
	for i, mu := range c.mus {
		mu.Lock()
		c.lrus[i].RemoveOldest()
		mu.Unlock()
	}
}

func (c *CacheShard) Len() int {
	l := 0
	for _, lru := range c.lrus {
		l += lru.Len()
	}
	return l
}

func (c *CacheShard) ShardLen() []int {
	l := make([]int, c.shardNum)
	for i, lru := range c.lrus {
		l[i] = lru.Len()
	}
	return l
}

func (c *CacheShard) Clear() {
	for i, mu := range c.mus {
		mu.Lock()
		c.lrus[i].Clear()
		mu.Unlock()
	}
}
