package data

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Red struct {
	Address  string      // The address to redis
	Password string      //  The password to redis
	Redis    *redis.Pool // This is our generic redis connection
}

func NewRedis(address string, password string) *Red {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
	}

	p := &Red{
		address,
		password,
		pool,
	}

	return p
}
