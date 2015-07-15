package database

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	groupLimit = 2
)

// TopNewsRankerResult data struct from mongodb aggregation result
// { "_id" : null, "total" : 326554 }
type TopNewsRankerResult struct {
	ID    string `bson:"_id"`
	Total int
	Items []map[string]bson.ObjectId
}

// TopNewsRanker main news top ranker
func TopNewsRanker() []string {
	fmt.Println("topnewsranker handled!")
	start := time.Now()

	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	gte := time.Now().Add(-time.Hour * hoursPerDayQuery)
	lte := time.Now()
	fmt.Println("query for this times gte: ", gte, " lte: ", lte)

	// query := []bson.M{
	// 	{"$match": bson.M{
	// 		"created_at": bson.M{"$gte": gte, "$lte": lte},
	// 	}},
	// 	{"$group": bson.M{
	// 		"_id": "$name",
	// 		"sum": bson.M{"$sum": 1},
	// 		"items": bson.M{
	// 			"$push": bson.M{
	// 				"temp":        "$main.temp",
	// 				"created_at":  "$created_at",
	// 				"description": "$weather"}},
	// 	}},
	// }
	// db.news_main.aggregate([{$group:{_id:"$category.initial"}}])
	// query for the meantime the latest of each group.
	// add to query for higher views
	var results []TopNewsRankerResult
	query := []bson.M{
		{"$match": bson.M{
			"created_at":       bson.M{"$gte": gte, "$lte": lte},
			"category.initial": bson.M{"$ne": " "},
			"score":            bson.M{"$gte": 1},
		}},
		{"$group": bson.M{
			"_id": "$category.initial",
			"sum": bson.M{"$sum": 1},
			"items": bson.M{
				"$push": bson.M{
					"_id": "$_id",
				},
			},
		}},
	}

	// pipe and execute the query
	c.Pipe(query).All(&results)
	fmt.Println("took: ", time.Since(start))
	return ExtractIDsFromResult(results...)
}

// ExtractIDsFromResult utils from extracting IDS from result
func ExtractIDsFromResult(results ...TopNewsRankerResult) []string {
	var extractIDs []string
	for _, i := range results {
		for iter, id := range i.Items {
			if iter >= groupLimit {
				break
			}
			extractIDs = append(extractIDs, id["_id"].Hex())
		}
	}
	return extractIDs
}
