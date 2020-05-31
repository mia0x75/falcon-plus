package g

import (
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// RedisConnPool Redis链接池对象
var RedisConnPool *redis.Pool

// InitRedisConnPool 初始化Redis链接池
func InitRedisConnPool() {
	cfg := Config().Redis

	addr := cfg.Addr
	maxIdle := cfg.MaxIdle
	waitTimeout := time.Duration(cfg.WaitTimeout) * time.Second

	RedisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   150,
		IdleTimeout: waitTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
}

// PingRedis 测试Redis链接
func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Errorf("[E] ping redis fail: %v", err)
	}
	return err
}
