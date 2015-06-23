package newsGetter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"web_apps/news_crawlers/modules/database"
)

var (
	hackerNewsProvider = "https://news.ycombinator.com"
	hackerNewsName     = "HackerNews"
)

// HackerNewsTopStoriesID struct for hacker news ids results
type HackerNewsTopStoriesID []int

// StartHackerNews starting GET hackernews
func StartHackerNews(loopCounterDelay int) {
	var wg sync.WaitGroup
	fmt.Println("starthacker news launched!")

	for t := range time.Tick(time.Duration(loopCounterDelay) * time.Second) {

		timeProfiler := make(chan string)

		topStoriesIds, err := topStoriesID()
		if err != nil {
			fmt.Println("skipping, err from topStoriesId")
			continue
		}
		fmt.Println("running the loop: ", t)

		c := make(chan int)
		for _, id := range topStoriesIds {
			wg.Add(1)
			go func(id int, timeProfiler chan string) {
				start := time.Now()
				newsContent := hackerNewsReader(id)
				ContentOutPut(newsContent, &wg)

				timeProfiler <- fmt.Sprintf("HN loop took: %v", time.Since(start))
			}(id, timeProfiler)
		}
		wg.Wait()
		close(c)
	}
}

// ContentOutPut data insert and db processing
func ContentOutPut(contentOutMsg jsonNewsBody, wg *sync.WaitGroup) {
	timeF := contentOutMsg.Time
	contentOutMsg.Time = int(time.Now().Unix())
	contentOutMsg.CreatedAt = fmt.Sprintf("%v", time.Now().Local())
	contentOutMsg.ProviderUrl = hackerNewsProvider
	contentOutMsg.ProviderName = hackerNewsName

	_ = timeF

	// check if can save
	// then save
	database.HackerNewsInsert(contentOutMsg, contentOutMsg.Title, wg)
}

// topStoriesId
func topStoriesID() ([]int, error) {
	var topStoriesIDURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	var idContainers HackerNewsTopStoriesID
	response, err := httpGet(topStoriesIDURL)
	if err != nil {
		var x []int
		return x, err
	}

	defer response.Body.Close()

	contents, _ := responseReader(response)
	if err := json.Unmarshal(contents, &idContainers); err != nil {
		return idContainers, nil
	}
	fmt.Printf("got %v ids:", len(idContainers))

	// make error handler
	return idContainers, nil
}

func hackerNewsReader(id int) jsonNewsBody {
	newsURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	var newsContent jsonNewsBody
	response, err := httpGet(newsURL)
	if err != nil {
		fmt.Println(err)
		return newsContent
	}
	defer response.Body.Close()

	contents, _ := responseReader(response)
	if err := json.Unmarshal(contents, &newsContent); err != nil {
		return newsContent
	}
	return newsContent
}
