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

	var wg sync.WaitGroup
	for t := range time.Tick(time.Duration(googleLoopCounterDelay) * time.Second) {
		_ = t

		c := make(chan int)
		fmt.Println("loop will start")
		for k, v := range TopicsList() {
			wg.Add(1)
			go func(k string, v TopicIdentity) {
				fmt.Println("running loop")
				GoogleNewsRequester(googleURLConstructor(v.Initial), v, &wg)
			}(k, v)
			// wg.Done()
		}
		wg.Wait()
		close(c)
	}
}

// GoogleNewsRequester google news http getter
func GoogleNewsRequester(url string, topic TopicIdentity, wsg *sync.WaitGroup) {
	defer wsg.Done()
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
	}

	GNResponse := googleNews.ResponseData
	var wg sync.WaitGroup
	for _, gn := range GNResponse.Results {
		// set news item category
		gn.Category = topic
		wg.Add(1)
		go func(gn GoogleNewsResults) {
			GoogleNewsDataSetter(gn, &wg)
		}(gn)
	}
	wg.Wait()
}

// GoogleNewsDataSetter builds and construct data for insertion
func GoogleNewsDataSetter(googleNews GoogleNewsResults, wg *sync.WaitGroup) {
	canSave := database.GoogleNewsFindIfExist(googleNews.Title)

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
	if canSave {
		saved := database.GoogleNewsInsert(jsonNews, wg)
		if saved {
			fmt.Println("saved!! google news!")
			return
		}
		fmt.Println("did not save!")
	}
}

//googleUrlConstructor return url string
func googleURLConstructor(v string) string {
	url := fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=%s&ned=jp&userip=192.168.0.1", v)
	return url
}
