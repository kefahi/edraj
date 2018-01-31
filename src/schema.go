package main

// SchemaMan to manage the various schemas in the system
type SchemaMan struct {
	mongoStore       MongoStore
	schemaCollection MongoCollection
}

func (sm *SchemaMan) init(mongoAddress string) {
	sm.mongoStore.init(mongoAddress, "schema")
	sm.schemaCollection.init("schema", &sm.mongoStore)
}
