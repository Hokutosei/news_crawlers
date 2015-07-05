package newsCache

import (
	"fmt"
	"strings"
	"web_apps/news_crawlers/modules/database"

	"gopkg.in/mgo.v2/bson"
)

var (
	newsIndexKeySlice = []string{"index", "ids"}
	newsIndexIDS      []bson.ObjectId
)

// NewsIndexCache make an index news for fast access
func NewsIndexCache(stop chan bool) {
	fmt.Println("starting news index cache...")

	result, err := database.NewsIndexNewsIDS()
	if err != nil {
		return
	}

	pushIDredis(result...)
}

func pushIDredis(IDS ...database.NewsIds) {
	conn := database.RedisPool.Get()
	defer conn.Close()

	key := RedisKeyGen(newsIndexKeySlice...)

	_, err := conn.Do("LPUSH", key, IDS)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// DeleteKey remove keys from redis
func DeleteKey(keys ...string) {
	conn := database.RedisPool.Get()
	defer conn.Close()

	for _, key := range keys {
		_, err := conn.Do("DEL", key)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// RedisKeyGen generate redis keys based on slice
func RedisKeyGen(keys ...string) string {
	return strings.Join(keys, ":")
}
