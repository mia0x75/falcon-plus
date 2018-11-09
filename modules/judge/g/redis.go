package g

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
)

var RedisConnPool *redis.Pool

func InitRedisConnPool() {
	cfg := Config().Redis
	addr := cfg.Addr
	maxIdle := cfg.MaxIdle
	waitTimeout := time.Duration(cfg.WaitTimeout) * time.Second
	connTimeout := time.Duration(cfg.ConnTimeout) * time.Second
	readTimeout := time.Duration(cfg.ReadTimeout) * time.Second
	writeTimeout := time.Duration(cfg.WriteTimeout) * time.Second

	RedisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: waitTimeout,
		Dial: func() (redis.Conn, error) {
			do := []redis.DialOption{
				redis.DialReadTimeout(readTimeout),
				redis.DialWriteTimeout(writeTimeout),
				redis.DialConnectTimeout(connTimeout),
			}
			c, err := redis.DialURL(addr, do...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
}

func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}
