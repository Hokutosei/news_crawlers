package newsCache

import (
	"fmt"
	"strings"
	"web_apps/news_crawlers/modules/database"
)

// NewsIndexCache make an index news for fast access
func NewsIndexCache(stop chan bool) {
	fmt.Println("starting news index cache...")

	result, err := database.NewsIndexNewsIDS()
	if err != nil {
		return
	}

	for _, i := range result {
		pushIDredis(i.ID.Hex())
	}
}

func pushIDredis(ID string) {
	conn := database.RedisPool.Get()
	defer conn.Close()

	keySlice := []string{"index", "ids"}
	key := RedisKeyGen(keySlice...)

	_, err := conn.Do("LPUSH", key, ID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// RedisKeyGen generate redis keys based on slice
func RedisKeyGen(keys ...string) string {
	return strings.Join(keys, ":")
}
