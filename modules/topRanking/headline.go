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
	daysAgoRetries  time.Duration = 5

	retries = 5
)

// Headlines get the most viewed/like news
func Headlines(loopDelay, loopRetry int) {
	fmt.Println("headline crawler started!")

	for t := range time.Tick(time.Duration(loopDelay) * time.Second) {
		fmt.Println(t)
		var idSlice []string
		var daysAgo = daysAgoDuration
		var daysTo time.Duration
		// searching := false

		// ensure we have enough idSlice len
		counter := 0
		for len(idSlice) < loopRetry {
			time.Sleep(time.Second * 5)
			idSlice = database.Headlines(daysAgo, daysTo)
			daysAgo++
			daysTo++
			fmt.Println(daysAgo)
			fmt.Println(daysAgoRetries)
			fmt.Println("retrying headline get... ", daysAgo)
			if counter >= retries {
				fmt.Println("retried so much headlines!")
				break
			}
			counter++
		}

		if len(idSlice) < 5 {
			return
		}

		key := newsCache.RedisKeyGen(headLines...)
		newsCache.PushIDredisObjectID(key, idSlice...)
	}
}
