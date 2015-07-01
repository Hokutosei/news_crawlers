package main

import (
	"fmt"
	"web_apps/news_crawlers/modules/config"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsGetter"
)

var (
	loopDelay = 300
)

// main entrypoint of the program
func main() {
	// strt getting analytics
	// InitNewRelic()

	config.StartEtcd()

	// connect to mongodb
	go database.MongodbStart()

	fmt.Println("starting!")
	// newsGetter.StartHackerNews(loopDelay)
	newsGetter.StartGoogleNews(loopDelay)
}
