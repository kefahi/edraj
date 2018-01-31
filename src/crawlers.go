package main

// CrawlersMan manage primary crawler and maintain the list of remote crawlers used
type CrawlersMan struct {
	mongoStore         MongoStore
	crawlersCollection MongoCollection
}

func (cm *CrawlersMan) init(mongoAddress string) {
	cm.mongoStore.init(mongoAddress, "crawlers")
	cm.crawlersCollection.init("crawlers", &cm.mongoStore)
}
