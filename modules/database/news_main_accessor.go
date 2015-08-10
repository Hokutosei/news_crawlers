package database

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"web_apps/news_crawlers/modules/utils"
)

var (

	// NewsMainCollection collection name
	NewsMainCollection               = "news_main"
	dayHours           time.Duration = 24
	hoursPerDayQuery                 = dayHours * 2
)

// NewsMainIndexNews responder for index news query
func NewsMainIndexNews() (AggregatedNews, error) {
	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	var aggregatedNews AggregatedNews
	// fix sorting query with
	// iter := coll.Find(nil).Sort(bson.D{{"field1", 1}, {"field2", -1}}).Iter()
	// refactor querying by including explicitly gte & lte
	// time.Local = time.UTC
	// gte := time.Now().Add(-time.Hour * hoursPerDayQuery)
	// lte := time.Now()
	// fmt.Println(gte)
	// fmt.Println(lte)
	// err := c.Find(bson.M{"url": bson.M{"$ne": ""}, "createdat": bson.M{"$gt": gte, "$lt": lte}}).Sort("-_id", "-score").Limit(searchLimitItems).All(&aggregatedNews)
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Sort("-_id", "-score").Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}

// NewsIds struct for result IDS only
type NewsIds struct {
	ID bson.ObjectId `bson:"_id"`
}

// NewsIndexNewsIDS retrieve all news index ids and save to cache MAIN
func NewsIndexNewsIDS() ([]NewsIds, error) {
	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	var aggregatedNews []NewsIds
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Select(bson.M{"_id": 1}).Sort("-_id", "-score", "image_url").Limit(searchLimitItems).All(&aggregatedNews)
	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}

//GetterNewsMainTopScore main top page news getter
func GetterNewsMainTopScore() (AggregatedNews, error) {
	c := MongodbSession.DB(Db).C(NewsMainCollection)
	var aggregatedNews AggregatedNews
	err := c.Find(bson.M{"url": bson.M{"$ne": ""}}).Sort("-score").Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}

//IncrementNewsScore increment news score
// increment news ite page view
func IncrementNewsScore(paramsID string) {
	c := MongodbSession.DB(Db).C(NewsMainCollection)
	var aggregatedNews interface{}

	err := c.Update(bson.M{"_id": bson.ObjectIdHex(paramsID)},
		bson.M{"$inc": bson.M{"score": 1}, "$currentDate": bson.M{"lastModified": true}})

	utils.HandleError(err)

	err = c.Find(bson.M{"_id": bson.ObjectIdHex(paramsID)}).One(&aggregatedNews)
	utils.HandleError(err)

	fmt.Println(aggregatedNews)
}

// NewsItemPage get news item data
func NewsItemPage(paramsID string) (interface{}, error) {
	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	var newsItem interface{}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(paramsID)}).One(&newsItem)

	if err != nil {
		return newsItem, err
	}

	return newsItem, nil
}

// GetCategorizedNews will get news with category news initials
func GetCategorizedNews(initial string) (AggregatedNews, error) {
	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	var aggregatedNews AggregatedNews
	err := c.Find(bson.M{"category.initial": initial}).Sort("-_id").Limit(searchLimitItems).All(&aggregatedNews)

	if err != nil {
		fmt.Println(err)
		return aggregatedNews, err
	}
	return aggregatedNews, nil
}
