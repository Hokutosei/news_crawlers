package newsGetter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"web_apps/news_crawlers/modules/database"
)

// GoogleNews interface for google news
type GoogleNews map[string]interface{}

// GoogleNewsResponseData response struct
type GoogleNewsResponseData struct {
	ResponseData struct {
		Results []GoogleNewsResults
	}
}

// ResponseData response struct
type ResponseData struct {
	Results []GoogleNewsResults
}

var (
	googleNewsProvider = "https://news.google.com/"
	googleNewsName     = "GoogleNews"
)

// TopicsList return a list of topics/categories
func TopicsList() Topics {
	topics := Topics{
		"society":       TopicIdentity{"y", "社会"},
		"international": TopicIdentity{"w", "国際"},
		"business":      TopicIdentity{"b", "ビジネス"},
		"politics":      TopicIdentity{"p", "政治"},
		"entertainment": TopicIdentity{"e", "エンタメ"},
		"sports":        TopicIdentity{"s", "スポーツ"},
		"technology":    TopicIdentity{"t", "テクノロジー"},
		"pickup":        TopicIdentity{"ir", "ピックアップ"},
	}

	return topics
}

// StartGoogleNews start collecting google news
func StartGoogleNews(googleLoopCounterDelay int) {
	fmt.Println("startgoogle news launched!")

	for t := range time.Tick(time.Duration(googleLoopCounterDelay) * time.Second) {
		_ = t

		fmt.Println("loop will start")
		var wsg sync.WaitGroup
		n := make(chan GoogleNewsResponseData)
		for _, v := range TopicsList() {
			wsg.Add(1)
			go func(v TopicIdentity) {
				go GoogleNewsRequester(googleURLConstructor(v.Initial), v, n, &wsg)

				result := <-n
				GoogleNewsRW(result, v)
			}(v)
		}
		wsg.Wait()
		close(n)
	}
}

// GoogleNewsRequester google news http getter
func GoogleNewsRequester(url string, topic TopicIdentity, c chan GoogleNewsResponseData, wg *sync.WaitGroup) {
	defer wg.Done()
	var googleNews GoogleNewsResponseData
	response, err := httpGet(url)
	if err != nil {
		fmt.Println("got error google!")
		return
	}

	defer response.Body.Close()

	contents, _ := responseReader(response)
	if err := json.Unmarshal(contents, &googleNews); err != nil {
		//return id_containers
		fmt.Println(err)
		c <- googleNews
		return
	}
	c <- googleNews
}

// GoogleNewsRW read and write data from google news
func GoogleNewsRW(gn GoogleNewsResponseData, topic TopicIdentity) {
	var wg sync.WaitGroup
	for _, g := range gn.ResponseData.Results {
		// set news item category
		g.Category = topic
		wg.Add(1)
		GoogleNewsDataSetter(g, &wg)
	}
	wg.Wait()
}

// GoogleNewsDataSetter builds and construct data for insertion
func GoogleNewsDataSetter(googleNews GoogleNewsResults, wg *sync.WaitGroup) {
	defer wg.Done()
	jsonNews := &jsonNewsBody{
		Title:          googleNews.Title,
		By:             "GoogleNews",
		Score:          0,
		Time:           int(time.Now().Unix()),
		Url:            googleNews.URL,
		ProviderName:   googleNewsName,
		RelatedStories: googleNews.RelatedStories,
		CreatedAt:      fmt.Sprintf("%v", time.Now().Local()),
		Category:       googleNews.Category,
		Image:          googleNews.Image,
	}

	// check if data exists already, need refactoring though
	saved := database.GoogleNewsInsert(jsonNews, googleNews.Title)
	if saved {
		fmt.Println("saved!! google news!")
		return
	}

	fmt.Println("did not save!")
}

//googleUrlConstructor return url string
func googleURLConstructor(v string) string {
	url := fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=%s&ned=jp&userip=192.168.0.1", v)
	return url
}
