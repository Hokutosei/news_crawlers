package database

import (
	"flag"
	"fmt"
	"time"
	_ "web_apps/news_crawlers/modules/config"

	"web_apps/news_crawlers/modules/utils"

	"github.com/garyburd/redigo/redis"
)

var (
	// RedisPool main redis pool connection
	RedisPool *redis.Pool

	redisServer   = flag.String("redisServer", ":6379", "")
	redisPassword = flag.String("redisPassword", "", "")
	redisHostKey  = "redisHost"
)

// StartRedis start connecting to redis
func StartRedis() {
	utils.Info(fmt.Sprintf("starting redis.."))
	redisHost := make(chan string)
	go GetRedisHost(redisHost)

	s := <-redisHost
	RedisPool = NewPool(s)
	utils.Info(fmt.Sprintf("connected to redis.."))
}

// NewPool create redis pool servers
func NewPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// GetRedisHost get redis host from etcd
func GetRedisHost(host chan string) {
	// redisHost, err := config.EtcdRawGetValue(redisHostKey)
	// if err != nil {
	// 	panic(err)
	// }

	host <- "130.211.255.234:6379"
}
