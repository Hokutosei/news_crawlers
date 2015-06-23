package main

import (
	"fmt"
	"web_apps/news_crawlers/modules/newsGetter"
)

var (
	loopDelay = 10
)

// main entrypoint of the program
func main() {

	fmt.Println("starting!")
	go newsGetter.StartHackerNews(loopDelay)
}
