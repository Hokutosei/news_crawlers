package database

import (
	"fmt"
	"sync"

	mongodb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GoogleNews interface for google news
type GoogleNews interface{}

var (
	googleNewsCollection = "news_main"
)

// GoogleNewsInsert insert data for google news
func GoogleNewsInsert(hn GoogleNews, URL string, wg *sync.WaitGroup) bool {
	sc := SessionCopy()
	c := sc.DB(Db).C(googleNewsCollection)
	defer sc.Close()

	// if !GoogleNewsFindIfExist(URL, sc) {
	// 	return false
	// }

	_, err := c.Upsert(bson.M{"url": URL}, hn)
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
func GoogleNewsFindIfExist(URL string, sc *mongodb.Session) bool {
	c := sc.DB(Db).C(googleNewsCollection)

	var result map[string]interface{}
	c.Find(bson.M{"url": URL}).One(&result)
	if result["url"] != nil && result["url"] != URL {
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
