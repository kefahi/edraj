package main

import (
	"path"
)

// AddonsMan to manage the addons installed on the system
type AddonsMan struct {
	mongoStore       MongoStore
	addonsCollection MongoCollection

	fileStore Storage
}

func (am *AddonsMan) init(mongoAddress string, rootPath string) {
	am.mongoStore.init(mongoAddress, "addons")

	am.addonsCollection = MongoCollection{}
	am.addonsCollection.init("crawlers", &am.mongoStore)

	am.fileStore.RootPath = path.Join(rootPath, "addons")
	am.fileStore.TrashPath = path.Join(rootPath, "trash", "addons")
}
