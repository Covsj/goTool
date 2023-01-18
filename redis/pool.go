package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var RedisPool *redis.Pool

func NewRedisPool(server string, port int, maxIdle, maxActive int) *redis.Pool {
	RedisPool = &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) {
			hostPort := fmt.Sprintf("%s:%d", server, port)
			return redis.Dial("tcp", hostPort, redis.DialConnectTimeout(time.Millisecond*150),
				redis.DialReadTimeout(time.Millisecond*150),
				redis.DialWriteTimeout(time.Millisecond*150))
		},
	}
	return RedisPool
}
