package main

// DomainsMan manage the list of known domains and the local domain
type DomainsMan struct {
	mongoStore        MongoStore
	domainsCollection MongoCollection
}

func (dm *DomainsMan) init(mongoAddress string) {
	dm.mongoStore.init(mongoAddress, "domains")
	dm.domainsCollection.init("domains", &dm.mongoStore)
}

// Server : Edraj server setup, can host multiple domain-legs
type Server struct {
	publicIPs  []string
	privateIPs []string
}
