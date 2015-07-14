package main

import (
	"fmt"
	"strconv"
	"web_apps/news_crawlers/modules/config"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsGetter"
)

var (
	loopDelay = 300
)

// CalcLoopDlay get from ENV our loop delay
// or from something like redis
func CalcLoopDlay() int {
	delay := config.GetEnvVar("loopDelay")
	i, err := strconv.Atoi(delay)
	if err != nil {
		return loopDelay
	}
	return i
}

// main entrypoint of the program
func main() {
	fmt.Println("starting!")
	// strt getting analytics
	// InitNewRelic()

	config.StartEtcd()

	// connect to mongodb
	go database.MongodbStart()
	go database.StartRedis()

	// newsGetter.StartHackerNews(loopDelay)

	newsGetter.StartGoogleNews(CalcLoopDlay())
}
