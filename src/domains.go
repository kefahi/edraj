package main

// DomainsMan manage the list of known domains and the local domain
// Server : Edraj server setup, can host multiple domain-legs
type DomainsMan struct {
	publicIPs         []string
	privateIPs        []string
	mongoStore        MongoStore
	domainsCollection MongoCollection
}

/*
func (dm *DomainsMan) init(mongoAddress string) {
	dm.mongoStore.init(mongoAddress, "domains")
	dm.domainsCollection.init("domains", &dm.mongoStore)
}*/

func (dm *DomainsMan) init(config *Config) (err error)                  { return }
func (dm *DomainsMan) query(request *Request) (response *QueryResponse) { return }
func (dm *DomainsMan) get(request *Request) (response *QueryResponse)   { return }
func (dm *DomainsMan) create(request *Request) (response Response)      { return }
func (dm *DomainsMan) update(request *Request) (response Response)      { return }
func (dm *DomainsMan) delete(request *Request) (response Response)      { return }
