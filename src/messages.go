package main

import (
	"path"
)

// MessagesMan all messaging goes here
type MessagesMan struct {
	mongoStore            MongoStore
	messagesCollection    MongoCollection
	attachmentsCollection MongoCollection
	fileStore             Storage
}

func (mm *MessagesMan) init(mongoAddress string, rootPath string) {
	mm.mongoStore.init(mongoAddress, "messages")
	mm.messagesCollection.init("messages", &mm.mongoStore)

	mm.fileStore.RootPath = path.Join(rootPath, "messags")
	mm.fileStore.TrashPath = path.Join(rootPath, "trash", "messages")
}
