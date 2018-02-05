package main

// MessagesMan all messaging goes here
type MessagesMan struct {
	mongoStore            MongoStore
	messagesCollection    MongoCollection
	attachmentsCollection MongoCollection
	fileStore             Storage
}

/*
func (mm *MessagesMan) init(mongoAddress string, rootPath string) {
	mm.mongoStore.init(mongoAddress, "messages")
	mm.messagesCollection.init("messages", &mm.mongoStore)

	mm.fileStore.RootPath = path.Join(rootPath, "messags")
	mm.fileStore.TrashPath = path.Join(rootPath, "trash", "messages")
}*/

func (mm *MessagesMan) init(config *Config) (err error)                  { return }
func (mm *MessagesMan) query(request *Request) (response *QueryResponse) { return }
func (mm *MessagesMan) get(request *Request) (response *QueryResponse)   { return }
func (mm *MessagesMan) create(request *Request) (response Response)      { return }
func (mm *MessagesMan) update(request *Request) (response Response)      { return }
func (mm *MessagesMan) delete(request *Request) (response Response)      { return }
