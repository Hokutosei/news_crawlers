package newsGetter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsCache"
)

// GoogleNews interface for google news
type GoogleNews map[string]interface{}

// GoogleNewsResponseData response struct
type GoogleNewsResponseData struct {
	Category     TopicIdentity
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
	fmt.Println(googleLoopCounterDelay)

	for t := range time.Tick(time.Duration(googleLoopCounterDelay) * time.Second) {
		_ = t
		fmt.Println("loop will start")

		var wsg sync.WaitGroup
		n := make(chan GoogleNewsResponseData)
		// cs := make(chan int)
		for _, v := range TopicsList() {
			wsg.Add(1)
			go func(v TopicIdentity) {
				go GoogleNewsRequester(googleURLConstructor(v.Initial), v, n, &wsg)

				result := <-n
				GoogleNewsRW(result, &wsg)
			}(v)
		}
		wsg.Wait()
		close(n)

		// cache index news keys
		newsCache.NewsIndexCache()
	}
}

// GoogleNewsRequester google news http getter
func GoogleNewsRequester(url string, topic TopicIdentity, c chan GoogleNewsResponseData, wg *sync.WaitGroup) {
	var googleNews GoogleNewsResponseData
	response, err := httpGet(url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("got error google!")
		wg.Done()
		return
	}

	defer response.Body.Close()
	contents, err := responseReader(response)

	if err != nil {
		fmt.Println(err)
		wg.Done()
		return
	}

	err = json.Unmarshal(contents, &googleNews)
	if err != nil {
		//return id_containers
		fmt.Println(err)
		// c <- googleNews
		wg.Done()
		return
	}

	// explicitly set gn news category
	googleNews.Category = topic
	c <- googleNews
}

// GoogleNewsRW read and write data from google news
func GoogleNewsRW(gn GoogleNewsResponseData, wg *sync.WaitGroup) {
	var wsg sync.WaitGroup
	for _, g := range gn.ResponseData.Results {
		// set news item category
		g.Category = gn.Category
		g.SecondaryTitle = g.Title
		wsg.Add(1)
		fmt.Println("-------- category ", gn.Category.Name)
		fmt.Println(g.Title)
		GoogleNewsDataSetter(g, &wsg)
	}
	wsg.Wait()
	wg.Done()
}

// GoogleNewsDataSetter builds and construct data for insertion
func GoogleNewsDataSetter(googleNews GoogleNewsResults, wg *sync.WaitGroup) {
	// defer wg.Done()
	start := time.Now()
	jsonNews := &jsonNewsBody{
		Title:          googleNews.Title,
		SecondaryTitle: googleNews.SecondaryTitle,
		By:             "GoogleNews",
		Score:          0,
		Time:           int(time.Now().Unix()),
		Url:            googleNews.URL,
		ProviderName:   googleNewsName,
		Publisher:      googleNews.Publisher,
		RelatedStories: googleNews.RelatedStories,
		CreatedAt:      fmt.Sprintf("%v", time.Now().Local()),
		// news refactoring of news Category from here
		Category: googleNews.Category,
		Image:    googleNews.Image,
	}

	// check if data exists already, need refactoring though
	saved := database.GoogleNewsInsert(jsonNews, googleNews.Title, wg)

	if saved {
		end := time.Since(start)
		fmt.Println("saved!! google news! took: ", end)
		// wg.Done()
		return
	}

	// wg.Done()
	fmt.Println("did not save!")
}

//googleUrlConstructor return url string
func googleURLConstructor(v string) string {
	// https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=t&ned=jp&userip=127.0.0.1
	url := fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=%s&ned=jp&userip=127.0.0.2", v)
	return url
}
