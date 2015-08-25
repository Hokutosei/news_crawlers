package database

import (
	"fmt"
	"math/rand"
	"news_worker/lib/utils"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var ()

// SuggestedRandItems struct holder
// { "_id" : null, "total" : 326554 }
type SuggestedRandItems struct {
	ID    string `bson:"_id"`
	Total int
	Items []map[string]bson.ObjectId
}

// SuggestRand query suggestion random news items
func SuggestRand(from time.Duration, to time.Duration) []string {
	fmt.Println("topnewsranker handled!")
	start := time.Now()

	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	fromVal := dayHours * from
	toVal := dayHours * to
	gte := time.Now().Add(-time.Hour * fromVal)
	lte := time.Now().Add(-time.Hour * toVal)
	fmt.Println("query for this times gte: ", gte, " lte: ", lte)

	var results []SuggestedRandItems

	// improve this count, kinda buggy
	collectionCount, err := c.Count()
	fmt.Println(collectionCount)
	if err != nil {
		fmt.Println(err)
		collectionCount = 2000
	}
	for i := 0; i < 10; i++ {
		var s SuggestedRandItems
		c.Find(bson.M{}).Sort("-_id").Skip(randomSkip(0, collectionCount)).One(&s)
		fmt.Println(s)
		results = append(results, s)
	}

	// pipe and execute the query
	// c.Pipe(query).All(&results)
	fmt.Println("took: ", time.Since(start))

	var extractIDs []string
	for _, i := range results {
		for iter, id := range i.Items {
			if iter >= groupLimit {
				break
			}
			extractIDs = append(extractIDs, id["_id"].Hex())
		}
	}
	fmt.Println(extractIDs)
	return extractIDs
}

func randomSkip(min, max int) int {
	rand.Seed(time.Now().Unix())
	value := rand.Intn(max-min) + min
	utils.Info(fmt.Sprintf("%v", value))
	return value
}
