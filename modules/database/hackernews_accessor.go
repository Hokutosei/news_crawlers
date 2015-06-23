package database

import (
	"fmt"
	"sync"

	mongodb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// HackerNews interface for hn news item
type HackerNews interface{}

// AggregatedNews interface struct for AggregatedNews
type AggregatedNews []interface{}

var (
	hackerNewsCollection = "news_main"
	searchLimitItems     = 50
)

// HackerNewsInsert insert data to mongodb
func HackerNewsInsert(hn HackerNews, hnTitle string, wg *sync.WaitGroup) {
	defer wg.Done()
	sc := SessionCopy()
	sc.SetMode(mongodb.Eventual, true)

	defer sc.Close()
	c := sc.DB(Db).C(hackerNewsCollection)

	if !HackerNewsFindIfExist(hnTitle, sc) {
		return
	}

	err := c.Insert(hn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("saved! hackernews!")
}

// HackerNewsFindIfExist check if data exists already before saving
func HackerNewsFindIfExist(title string, sc *mongodb.Session) bool {
	c := sc.DB(Db).C(hackerNewsCollection)

	var result map[string]interface{}
	c.Find(bson.M{"title": title}).One(&result)
	if result["title"] != nil {
		return false
	}
	return true
}

// HackerNewsIndexNews index page data
func HackerNewsIndexNews() (AggregatedNews, error) {
	sc := SessionCopy()
	defer sc.Close()
	c := sc.DB(Db).C(hackerNewsCollection)

	var aggregatedNews AggregatedNews
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Sort("-score").Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}

// HackerNewsLatestNews aggregate latest news
func HackerNewsLatestNews() (AggregatedNews, error) {
	sc := SessionCopy()
	defer sc.Close()
	c := sc.DB(Db).C(hackerNewsCollection)

	var aggregatedNews AggregatedNews
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Sort("-_id").Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}

// HackerNewsFeedMore feed more data on index list
func HackerNewsFeedMore(contentType string, length int) (AggregatedNews, error) {
	q := map[string]string{"latest_news": "-_id", "top_score_news": "-score"}
	sc := SessionCopy()
	defer sc.Close()
	c := sc.DB(Db).C(hackerNewsCollection)

	var aggregatedNews AggregatedNews
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Skip(length).Sort(q[contentType]).Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}
