package topRanking

import (
	"fmt"
	"time"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsCache"
)

var (
	// main todayTopRank news key
	todayTopRank = []string{"index", "news_top_rank"}
)

// GenerateTopRanking aggregate all news categories
// to index by grouping them and limiting
func GenerateTopRanking(loopDelay int) {
	fmt.Println("GenerateTopRanking starting...")

	for t := range time.Tick(time.Duration(loopDelay) * time.Second) {
		fmt.Println(t)
		idSlice := database.TopNewsRanker()

		key := newsCache.RedisKeyGen(todayTopRank...)
		newsCache.PushIDredisObjectID(key, idSlice...)
	}
}
