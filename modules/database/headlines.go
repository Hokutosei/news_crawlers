package database

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// HeadlineNewsStruct headline news item struct
type HeadlineNewsStruct struct {
	ID    string `bson:"_id"`
	Items []map[string]bson.ObjectId
	// Title string `bson:"title" json:"title"`
	// Image struct {
	// 	URL string `bson:"url"`
	// } `bson:"image"`
}

// Headlines index top headline news
func Headlines(from time.Duration, to time.Duration) []string {
	fmt.Println("headline crawler handled!")
	start := time.Now()

	sc := SessionCopy()
	c := sc.DB(Db).C(NewsMainCollection)
	defer sc.Close()

	fromVal := dayHours * from
	toVal := dayHours * to
	gte := time.Now().Add(-time.Hour * fromVal)
	lte := time.Now().Add(-time.Hour * toVal)
	fmt.Println("query for this times gte: ", gte, " lte: ", lte)

	// db.news_main.aggregate([{$match:{"created_at": {$gte: x, $lte: e}, "image.url": {$ne: ""}, "image.tbheight":{$gte:50}}}])
	var results []HeadlineNewsStruct
	query := []bson.M{
		{"$match": bson.M{
			// "created_at": bson.M{"$gte": gte, "$lte": lte},
			"image.url": bson.M{"$ne": ""},
			// "image.tbheight": bson.M{"$gte": 50},
			"score": bson.M{"$gte": 1},
		}},
		{"$group": bson.M{
			"_id": "$category.initial",
			"sum": bson.M{"$sum": 1},
			"items": bson.M{
				"$push": bson.M{
					"_id": "$_id",
					// "score": "$score",
				},
			},
		}},
		// {"$sort": bson.M{
		// 	"score": -1,
		// }},
	}

	// pipe and execute the query
	c.Pipe(query).All(&results)
	fmt.Println("has items", len(results), "took: ", time.Since(start))
	// return ExtractIDsFromResult(results...
	return extractID(results...)
}

func extractID(results ...HeadlineNewsStruct) []string {
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
