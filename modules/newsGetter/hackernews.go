package newsGetter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
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

		n := make(chan jsonNewsBody)
		failCounter := make(chan int)
		for _, id := range topStoriesIds {
			wg.Add(1)
			go func(id int, timeProfiler chan string) {
				hackerNewsReader(id, n, failCounter)
			}(id, timeProfiler)
			select {
			case newsContent := <-n:
				fmt.Println(1)
				ContentOutPut(newsContent, &wg)
			case fail := <-failCounter:
				fmt.Println(2)
				fmt.Println("failure: ", fail)
			default:
				fmt.Println(3)
				fmt.Println("news getter failed")
			}
			wg.Done()

		}
		wg.Wait()
		close(n)
		close(failCounter)
	}
}

// ContentOutPut data insert and db processing
func ContentOutPut(contentOutMsg jsonNewsBody, wg *sync.WaitGroup) {
	defer wg.Done()
	timeF := contentOutMsg.Time
	contentOutMsg.Time = int(time.Now().Unix())
	contentOutMsg.CreatedAt = fmt.Sprintf("%v", time.Now().Local())
	contentOutMsg.ProviderUrl = hackerNewsProvider
	contentOutMsg.ProviderName = hackerNewsName

	_ = timeF

	// check if can save
	// then save
	// database.HackerNewsInsert(contentOutMsg, contentOutMsg.Title, wg)
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

// hackerNewsReader http request to hn and read write to channel
func hackerNewsReader(id int, n chan jsonNewsBody, fail chan int) {
	newsURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	var newsContent jsonNewsBody
	response, err := httpGet(newsURL)
	if err != nil {
		// fmt.Println(err)
		fail <- 1 + <-fail
		// n <- newsContent
		return
	}

	defer response.Body.Close()

	contents, _ := responseReader(response)
	if err := json.Unmarshal(contents, &newsContent); err != nil {
		// n <- newsContent
		fail <- 1 + <-fail
		return
	}
	n <- newsContent
}
