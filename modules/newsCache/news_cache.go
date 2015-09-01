package newsCache

import (
	"fmt"
	"strings"
	"time"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/utils"

	"gopkg.in/mgo.v2/bson"
)

var (
	newsIndexIDS []bson.ObjectId
)

// NewsIndexCache make an index news for fast access
func NewsIndexCache() {
	start := time.Now()
	fmt.Println("starting news index cache...")

	for _, lang := range database.Languages() {
		newsIndexKeySlice := []string{"index", "ids"}
		result, err := database.NewsIndexNewsIDS(lang)
		if err != nil {
			continue
		}
		newsIndexKeySlice = append(newsIndexKeySlice, lang)

		key := RedisKeyGen(newsIndexKeySlice...)
		PushIDredis(key, result...)
		utils.Info(fmt.Sprintf("news index language: %v cache took: %v", lang, time.Since(start)))

	}
}

// PushIDredis util to push bulk id to redis w/key
// TODO refactor this!
func PushIDredis(key string, IDS ...database.NewsIds) {
	start := time.Now()
	conn := database.RedisPool.Get()
	defer conn.Close()

	var strID []string
	for _, item := range IDS {
		strID = append(strID, item.ID.Hex())
	}
	// DELETE existing
	conn.Send("DEL", key)

	reversedIDs := ReverseSlice(strID...)
	for _, id := range reversedIDs {
		conn.Send("LPUSH", key, id)
	}
	conn.Flush()
	res, err := conn.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.Info(fmt.Sprintf("key: %v ; push to cache index took: %v and redis: %v", key, time.Since(start), res))
}

// PushIDredisObjectID util to push bulk id to redis w/key
// TODO refactor this!
func PushIDredisObjectID(key string, IDs ...string) {
	start := time.Now()
	conn := database.RedisPool.Get()
	defer conn.Close()

	if len(IDs) == 0 {
		return
	}

	// DELETE existing
	conn.Send("DEL", key)

	reversedIDs := ReverseSlice(IDs...)
	for _, id := range reversedIDs {
		conn.Send("LPUSH", key, id)
	}
	conn.Flush()
	res, err := conn.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("push to cache index took: ", time.Since(start), "and redis: ", res)
	utils.Info(fmt.Sprintf("key: [%v] ; push to cache index took: %v and redis: %v", key, time.Since(start), res))
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
	return s
}

// RedisKeyGen generate redis keys based on slice
func RedisKeyGen(keys ...string) string {
	return strings.Join(keys, ":")
}
