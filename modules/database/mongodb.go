package database

import (
	"fmt"

	mongodb "gopkg.in/mgo.v2"

	"web_apps/news_crawlers/modules/config"
)

var (
	// MongodbSession main mongodb connection session
	MongodbSession *mongodb.Session

	// Db database name
	Db = "news_aggregator"

	//mongodbClusterKey etcd key name
	mongodbClusterKey = "mongodb_cluster1"
)

// GetMongodbCluster retrieve mongodb cluster node from etcd
func GetMongodbCluster(host chan string) {
	mongodbCluster, err := config.EtcdRawGetValue(mongodbClusterKey)
	if err != nil {
		panic(err)
	}

	host <- mongodbCluster
}

// MongodbStart start connecting to mongodb
func MongodbStart() {
	fmt.Println("starting mongodb..")

	mongodbCluster := make(chan string)
	go GetMongodbCluster(mongodbCluster)

	host := <-mongodbCluster
	connectionStr := fmt.Sprintf("mongodb://%v/?maxPoolSize=10", host)
	// fmt.Println(connectionStr)
	// mongoDBDialInfo := &mongodb.DialInfo{
	// 	Addrs:   []string{connectionStr},
	// 	Timeout: 10 * time.Second,
	// }

	session, err := mongodb.Dial(connectionStr)
	if err != nil {
		// fmt.Println(err)
		panic(err)
	}
	fmt.Println("connected to mongodb...")
	MongodbSession = session
	// MongodbSession.SetMode(mongodb.Eventual, true)
}

// SessionCopy make copy of a mongodb session
func SessionCopy() *mongodb.Session {
	return MongodbSession.Copy()
}
