package main

import (
	"fmt"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsGetter"
)

var (
	loopDelay = 10
)

// main entrypoint of the program
func main() {
	// strt getting analytics
	InitNewRelic()
  
	// connect to mongodb
	go database.MongodbStart()

	fmt.Println("starting!")
	newsGetter.StartHackerNews(loopDelay)
}
