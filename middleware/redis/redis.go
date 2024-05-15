package redis

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

func Ping(pool *redis.Pool) (string, error) {
	conn := pool.Get()
	defer func() {
		conn.Close()
	}()
	return redis.String(conn.Do("ping"))
}

func SetEx(key, value interface{}, expireTime int) (err error) {
	conn := redisPool.Get()
	if conn == nil {
		return errors.New("get nil redis connection")
	}
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	retryTime := 2
	for i := 0; i < retryTime; i++ {
		if _, err = redis.String(conn.Do("SETEX", key, expireTime, value)); err != nil {
			if i == retryTime-1 {
				return err
			}
			time.Sleep(time.Duration(20) * time.Millisecond) // 等待20ms重试
			continue
		} else {
			break
		}
	}
	return nil
}
