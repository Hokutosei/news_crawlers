package topRanking

import (
	"fmt"
	"time"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsCache"
)

var (
	headLines                     = []string{"weekly", "headlines"}
	daysAgoDuration time.Duration = 7
	daysAgoRetries  time.Duration = 10
)

// Headlines get the most viewed/like news
func Headlines(loopDelay int) {
	for t := range time.Tick(time.Duration(loopDelay) * time.Second) {
		fmt.Println(t)
		var idSlice []string
		var daysAgo = daysAgoDuration
		var daysTo time.Duration
		// searching := false

		// ensure we have enough idSlice len
		for len(idSlice) < 5 {
			time.Sleep(time.Second * 5)
			idSlice = database.Headlines(daysAgo, daysTo)
			daysAgo++
			daysTo++
			if daysAgo >= daysAgoRetries {
				panic("retried so much")
			}
		}

		key := newsCache.RedisKeyGen(headLines...)
		newsCache.PushIDredisObjectID(key, idSlice...)
	}
}
