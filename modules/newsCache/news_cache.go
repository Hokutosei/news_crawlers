package newsCache

import (
	"fmt"
	"strings"
	"time"
	"web_apps/news_crawlers/modules/database"

	"gopkg.in/mgo.v2/bson"
)

var (
	newsIndexKeySlice = []string{"index", "ids"}
	newsIndexIDS      []bson.ObjectId
)

// NewsIndexCache make an index news for fast access
func NewsIndexCache() {
	fmt.Println("starting news index cache...")

	result, err := database.NewsIndexNewsIDS()
	if err != nil {
		return
	}
	pushIDredis(result...)
	// stop <- 0
}

func pushIDredis(IDS ...database.NewsIds) {
	start := time.Now()
	conn := database.RedisPool.Get()
	defer conn.Close()

	key := RedisKeyGen(newsIndexKeySlice...)
	var strID []string
	for _, item := range IDS {
		strID = append(strID, item.ID.Hex())
	}
	// DELETE existing
	conn.Send("DEL", key)

	reversedIDs := ReverseSlice(strID...)
	for _, id := range reversedIDs {
		fmt.Println(id)
		conn.Send("LPUSH", key, id)
	}
	conn.Flush()
	res, err := conn.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("push to cache index took: ", time.Since(start), "and redis: ", res)
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

// ReverseSlice util to reverse slice
func ReverseSlice(s ...string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	fmt.Println(s)
	return s
}

// RedisKeyGen generate redis keys based on slice
func RedisKeyGen(keys ...string) string {
	return strings.Join(keys, ":")
}
