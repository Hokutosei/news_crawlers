package database

import (
	"fmt"
	"math/rand"
	"time"
	"web_apps/news_crawlers/modules/utils"

	"gopkg.in/mgo.v2/bson"
)

var ()

// SuggestedRandItems struct holder
// { "_id" : null, "total" : 326554 }
type SuggestedRandItems struct {
	ID bson.ObjectId `bson:"_id"`
}

// SuggestRand query suggestion random news items
func SuggestRand(from time.Duration, to time.Duration) []string {
	start := time.Now()

	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	gte, lte := TimeRange(from, to)
	fmt.Println("query for this times gte: ", gte, " lte: ", lte)

	var results []SuggestedRandItems

	// improve this count, kinda buggy
	collectionCount, err := c.Find(bson.M{"created_at": bson.M{"$gte": gte, "$lte": lte}}).Count()
	if err != nil {
		fmt.Println(err)
		collectionCount = 2000
	}
	for i := 0; i < 10; i++ {
		var s SuggestedRandItems
		skipVal := randomSkip(0, collectionCount)
		c.Find(bson.M{"created_at": bson.M{"$gte": gte, "$lte": lte}}).Select(bson.M{"_id": 1}).Sort("-_id").Skip(skipVal).One(&s)
		results = append(results, s)
	}

	utils.Info(fmt.Sprintf("suggestRand took: %v", time.Since(start)))

	var extractIDs []string
	for _, item := range results {
		extractIDs = append(extractIDs, item.ID.Hex())
	}
	return extractIDs
}

func randomSkip(min, max int) int {
	value := min + rand.Intn(max-min)
	// utils.Info(fmt.Sprintf("%v", value))
	return value
}
