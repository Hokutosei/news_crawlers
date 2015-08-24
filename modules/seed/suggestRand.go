package seed

import (
	"fmt"
	"time"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsCache"
	"web_apps/news_crawlers/modules/utils"
)

var (
	suggestRandKey = []string{"suggest", "random"}
)

// SuggestRand suggest random items
func SuggestRand(loopDelay, retry int) {
	utils.Info(fmt.Sprintf("suggest random start"))

	for t := range time.Tick(time.Duration(loopDelay) * time.Second) {
		fmt.Println(t)
		var idSlice []string
		var daysAgo time.Duration = 1
		var daysTo time.Duration

		// ensure we have enough idSlice len
		for len(idSlice) < retry {
			time.Sleep(time.Second * 5)
			idSlice = database.SuggestRand(daysAgo, daysTo)
			daysAgo++
			daysTo++
			if daysAgo >= 5 {
				fmt.Println("retried so much top ranking")
				break
			}
		}

		if len(idSlice) < 5 {
			return
		}

		key := newsCache.RedisKeyGen(suggestRandKey...)
		newsCache.PushIDredisObjectID(key, idSlice...)
	}
}
