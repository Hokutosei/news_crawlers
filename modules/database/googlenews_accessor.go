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
func GoogleNewsInsert(hn GoogleNews, title string, imgURL string) bool {
	sc := SessionCopy()
	c := sc.DB(Db).C(googleNewsCollection)
	defer sc.Close()

	if !GoogleNewsFindIfExist(title, imgURL, sc) {
		return false
	}

	err := c.Insert(hn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// GoogleNewsFindIfExist check google news current data if exist before insert
func GoogleNewsFindIfExist(title string, imgURL string, sc *mongodb.Session) bool {
	c := sc.DB(Db).C(googleNewsCollection)

	var result map[string]interface{}
	// encodedTitle := utils.ToUtf8(title)

	query := bson.M{"title": title}
	if imgURL != "" {
		query = bson.M{"image_url": imgURL}
	}

	c.Find(query).One(&result)
	// c.Find(bson.M{"image_url": imgURL}).One(&result)

	// debug
	fmt.Println(result["image_url"])
	fmt.Println(result["title"])
	fmt.Println(title)
	fmt.Println(imgURL)
	fmt.Println(result == nil)
	fmt.Println("-----------------------------------------")
	// validate if any record found

	// return true if search is nil
	if result == nil {
		return true
	}

	// else do some check
	if result != nil ||
		title == result["title"] {
		return false
	}

	if result != nil ||
		result["image_url"] != " " ||
		result["image_url"] != nil ||
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
