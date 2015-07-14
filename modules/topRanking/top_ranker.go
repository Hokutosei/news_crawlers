package topRanking

import (
	"fmt"
	"time"
	"web_apps/news_crawlers/modules/database"
)

// GenerateTopRanking aggregate all news categories
// to index by grouping them and limiting
func GenerateTopRanking(loopDelay int) {
	fmt.Println("GenerateTopRanking starting...")

	for t := range time.Tick(time.Duration(loopDelay) * time.Second) {
		fmt.Println(t)
		database.TopNewsRanker()
	}
}
