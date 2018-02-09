package main

import (
	"path"

	mgo "gopkg.in/mgo.v2"
)

// MessagesMan all messaging goes here
type MessagesMan struct {
	mongoDb *mgo.Database
	//messagesCollection    MongoCollection
	//attachmentsCollection MongoCollection
	fileStore Storage
}

func (mm *MessagesMan) init(config *Config) (err error) {

	mm.mongoDb = mongoSession.DB(message)
	mm.fileStore.RootPath = path.Join(config.dataPath, message)
	mm.fileStore.TrashPath = path.Join(config.dataPath, trash, message)
	return
}
func (mm *MessagesMan) query(request *Request) (response *QueryResponse) { return }
func (mm *MessagesMan) get(request *Request) (response *QueryResponse)   { return }
func (mm *MessagesMan) create(request *Request) (response Response)      { return }
func (mm *MessagesMan) update(request *Request) (response Response)      { return }
func (mm *MessagesMan) delete(request *Request) (response Response)      { return }
