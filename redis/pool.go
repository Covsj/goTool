package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var RedisPool *redis.Pool

func NewRedisPool(server string, port int, maxIdle, maxActive, db int) *redis.Pool {
	RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(5) * time.Second,
		Dial: func() (redis.Conn, error) {
			hostPort := fmt.Sprintf("%s:%d", server, port)
			return redis.Dial("tcp", hostPort, redis.DialDatabase(db))
		},
	}
	return RedisPool
}
