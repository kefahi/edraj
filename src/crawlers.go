package main

// CrawlersMan manage primary crawler and maintain the list of remote crawlers used
type CrawlersMan struct {
	mongoStore         MongoStore
	crawlersCollection MongoCollection
}

/*
func (cm *CrawlersMan) init(mongoAddress string) {
	cm.mongoStore.init(mongoAddress, "crawlers")
	cm.crawlersCollection.init("crawlers", &cm.mongoStore)
}*/

func (cm *CrawlersMan) init(config *Config) (err error)                  { return }
func (cm *CrawlersMan) query(request *Request) (response *QueryResponse) { return }
func (cm *CrawlersMan) get(request *Request) (response *QueryResponse)   { return }
func (cm *CrawlersMan) create(request *Request) (response Response)      { return }
func (cm *CrawlersMan) update(request *Request) (response Response)      { return }
func (cm *CrawlersMan) delete(request *Request) (response Response)      { return }
