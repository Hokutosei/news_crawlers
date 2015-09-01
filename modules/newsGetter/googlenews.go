package newsGetter

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"web_apps/news_crawlers/modules/database"
	"web_apps/news_crawlers/modules/newsCache"
	"web_apps/news_crawlers/modules/utils"
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

// GetterLanguages language for news getters
func GetterLanguages() []string {
	//ned=en_ph
	return database.Languages()
}

// StartGoogleNews start collecting google news
func StartGoogleNews(googleLoopCounterDelay int) {
	utils.Info(fmt.Sprintf("startgoogle news launched!"))
	fmt.Println(googleLoopCounterDelay)

	for t := range time.Tick(time.Duration(googleLoopCounterDelay) * time.Second) {
		_ = t
		start := time.Now()
		utils.Info(fmt.Sprintf("google news loop start"))

		var wsg sync.WaitGroup
		n := make(chan GoogleNewsResponseData)
		languages := GetterLanguages()

		for _, v := range TopicsList() {
			wsg.Add(len(languages))
			go func(v TopicIdentity) {
				for _, lang := range languages {
					result, err := GoogleNewsGetter(googleURLConstructor(v.Initial, lang), v)
					if err != nil {
						wsg.Done()
						continue
					}
					GoogleNewsRW(result, &wsg, lang)
				}
			}(v)
		}
		wsg.Wait()
		close(n)

		// cache index news keys
		newsCache.NewsIndexCache()
		utils.Info(fmt.Sprintf("Google news took: %v", time.Since(start)))
	}
}

// GoogleNewsGetter request google api
func GoogleNewsGetter(url string, topic TopicIdentity) (GoogleNewsResponseData, error) {
	var googleNews GoogleNewsResponseData
	response, err := httpGet(url)
	if err != nil {
		utils.Info(fmt.Sprintf("error in google news get %v", err))
		return googleNews, err
	}
	defer response.Body.Close()

	contents, err := responseReader(response)
	if err != nil {
		utils.Info(fmt.Sprintf("error in response reader %v", err))
		return googleNews, err
	}

	err = json.Unmarshal(contents, &googleNews)
	if err != nil {
		utils.Info(fmt.Sprintf("error in unmarshal %v", err))
		return googleNews, err
	}

	googleNews.Category = topic
	return googleNews, nil
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
func GoogleNewsRW(gn GoogleNewsResponseData, wg *sync.WaitGroup, lang string) {
	for _, g := range gn.ResponseData.Results {
		// set news item category
		g.Category = gn.Category
		g.SecondaryTitle = g.Title
		g.Language = lang
		GoogleNewsDataSetter(g)
	}
	defer wg.Done()
}

// GoogleNewsDataSetter builds and construct data for insertion
func GoogleNewsDataSetter(googleNews GoogleNewsResults) {

	// main data struct for NEWS
	start := time.Now()
	jsonNews := &jsonNewsBody{
		Title:          googleNews.Title,
		SecondaryTitle: googleNews.SecondaryTitle,
		EncodedTitle:   utils.ToUtf8(googleNews.Title),
		By:             "GoogleNews",
		Score:          0,
		Time:           int(time.Now().Unix()),
		Url:            googleNews.URL,
		ImageURL:       googleNews.Image.URL,
		ProviderName:   googleNewsName,
		Publisher:      googleNews.Publisher,
		RelatedStories: googleNews.RelatedStories,
		CreatedAt:      time.Now().Local(),
		Category:       googleNews.Category,
		Image:          googleNews.Image,
		Content:        googleNews.Content,
		Language:       googleNews.Language,
	}

	// check if data exists already, need refactoring though
	saved := database.GoogleNewsInsert(jsonNews, googleNews.Title, googleNews.Image.URL)

	if saved {
		end := time.Since(start)
		fmt.Println("saved!! google news! took: ", end)
		return
	}

	fmt.Println("did not save!")
}

//googleUrlConstructor return url string
func googleURLConstructor(v, lang string) string {
	// https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=t&ned=jp&userip=127.0.0.1
	// url := fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=%s&ned=jp&userip=127.0.0.2", v)
	url := fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/news?v=1.0&topic=%s&ned=%v&userip=127.0.0.2", v, lang)
	return url
}
