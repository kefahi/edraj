package main

import mgo "gopkg.in/mgo.v2"

// DomainsMan manage the list of known domains and the local domain
// Server : Edraj server setup, can host multiple domain-legs
type DomainsMan struct {
	publicIPs  []string
	privateIPs []string
	mongoDb    *mgo.Database
}

func (dm *DomainsMan) init(config *Config) (err error) {
	dm.mongoDb = mongoSession.DB(domain)
	return
}
func (dm *DomainsMan) query(request *Request) (response *QueryResponse) { return }
func (dm *DomainsMan) get(request *Request) (response *QueryResponse)   { return }
func (dm *DomainsMan) create(request *Request) (response Response)      { return }
func (dm *DomainsMan) update(request *Request) (response Response)      { return }
func (dm *DomainsMan) delete(request *Request) (response Response)      { return }
