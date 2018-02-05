package main

import "path"

// MessagesMan all messaging goes here
type MessagesMan struct {
	mongoStore MongoStore
	//messagesCollection    MongoCollection
	//attachmentsCollection MongoCollection
	fileStore Storage
}

func (mm *MessagesMan) init(config *Config) (err error) {
	mm.mongoStore.init(config.mongoAddress, message)

	mm.fileStore.RootPath = path.Join(config.dataPath, message)
	mm.fileStore.TrashPath = path.Join(config.dataPath, trash, message)
	return
}
func (mm *MessagesMan) query(request *Request) (response *QueryResponse) { return }
func (mm *MessagesMan) get(request *Request) (response *QueryResponse)   { return }
func (mm *MessagesMan) create(request *Request) (response Response)      { return }
func (mm *MessagesMan) update(request *Request) (response Response)      { return }
func (mm *MessagesMan) delete(request *Request) (response Response)      { return }
