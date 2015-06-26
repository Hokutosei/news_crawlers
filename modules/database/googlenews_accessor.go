package database

import (
	"fmt"

	mongodb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GoogleNews interface for google news
type GoogleNews interface{}

var (
	googleNewsCollection = "news_main"
)

// GoogleNewsInsert insert data for google news
func GoogleNewsInsert(hn GoogleNews, title string) bool {
	sc := SessionCopy()
	c := sc.DB(Db).C(googleNewsCollection)
	defer sc.Close()

	if !GoogleNewsFindIfExist(title, sc) {
		return false
	}

	err := c.Insert(hn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
	//	fmt.Println("saved!")
}

// GoogleNewsFindIfExist check google news current data if exist before insert
func GoogleNewsFindIfExist(title string, sc *mongodb.Session) bool {
	c := sc.DB(Db).C(googleNewsCollection)

	var result map[string]interface{}
	c.Find(bson.M{"title": title}).One(&result)
	if result["title"] != nil {
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
