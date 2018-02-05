package main

// DomainsMan manage the list of known domains and the local domain
// Server : Edraj server setup, can host multiple domain-legs
type DomainsMan struct {
	publicIPs         []string
	privateIPs        []string
	mongoStore        MongoStore
	domainsCollection MongoCollection
}

func (dm *DomainsMan) init(mongoAddress string) {
	dm.mongoStore.init(mongoAddress, "domains")
	dm.domainsCollection.init("domains", &dm.mongoStore)
}
