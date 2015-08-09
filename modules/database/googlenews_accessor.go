package database

import (
	"fmt"
	"sync"
	"web_apps/news_crawlers/modules/utils"

	mongodb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GoogleNews interface for google news
type GoogleNews interface{}

var (
	googleNewsCollection = "news_main"
)

// GoogleNewsInsert insert data for google news
func GoogleNewsInsert(hn GoogleNews, title string, imgURL string, wg *sync.WaitGroup) bool {
	sc := SessionCopy()
	c := sc.DB(Db).C(googleNewsCollection)
	defer sc.Close()

	if !GoogleNewsFindIfExist(title, imgURL, sc) {
		wg.Done()
		return false
	}

	err := c.Insert(hn)
	if err != nil {
		fmt.Println(err)
		wg.Done()
		return false
	}
	wg.Done()
	return true
	//	fmt.Println("saved!")
}

// GoogleNewsFindIfExist check google news current data if exist before insert
func GoogleNewsFindIfExist(title string, imgURL string, sc *mongodb.Session) bool {
	c := sc.DB(Db).C(googleNewsCollection)

	var result map[string]interface{}
	encodedTitle := utils.ToUtf8(title)
	c.Find(bson.M{"encoded_title": encodedTitle}).One(&result)
	// c.Find(bson.M{"image_url": imgURL}).One(&result)

	// debug
	fmt.Println(imgURL)
	fmt.Println(result["image_url"])
	fmt.Println(title)
	fmt.Println(result["title"])
	fmt.Println("-----------------------------------------")

	// validate if any record found
	if title == result["title"] ||
		len(result) != 0 ||
		imgURL == result["image_url"] {
		return false
	}

	return true
}

// GoogleNewsIndexNews aggregated news list for google news
func GoogleNewsIndexNews() (AggregatedNews, error) {
	sc := SessionCopy()
	c := sc.DB(Db).C(googleNewsCollection)
	defer sc.Close()

	var aggregatedNews AggregatedNews
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Sort("-score").Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}
