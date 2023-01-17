package redis

import "github.com/garyburd/redigo/redis"

func Ping(pool *redis.Pool) (string, error) {
	conn := pool.Get()
	defer func() {
		conn.Close()
	}()
	return redis.String(conn.Do("ping"))
}
