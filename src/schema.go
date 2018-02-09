package main

import mgo "gopkg.in/mgo.v2"

// SchemaMan to manage the various schemas in the system
type SchemaMan struct {
	mongoDb *mgo.Database
}

/*
func (sm *SchemaMan) init(mongoAddress string) {
	sm.mongoStore.init(mongoAddress, "schema")
	sm.schemaCollection.init("schema", &sm.mongoStore)
}*/

func (sm *SchemaMan) init(config *Config) (err error) {
	sm.mongoDb = mongoSession.DB(schema)
	return
}

func (sm *SchemaMan) query(request *Request) (response *QueryResponse) { return }
func (sm *SchemaMan) get(request *Request) (response *QueryResponse)   { return }
func (sm *SchemaMan) create(request *Request) (response Response)      { return }
func (sm *SchemaMan) update(request *Request) (response Response)      { return }
func (sm *SchemaMan) delete(request *Request) (response Response)      { return }
