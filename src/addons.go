package main

import "path"

// AddonsMan to manage the addons installed on the system
type AddonsMan struct {
	mongoStore       MongoStore
	addonsCollection MongoCollection

	fileStore Storage
}

/*
func (am *AddonsMan) init(mongoAddress string, rootPath string) {
}
*/

func (am *AddonsMan) init(config *Config) (err error) {
	am.mongoStore.init(config.mongoAddress, addon)

	am.addonsCollection = MongoCollection{}
	am.addonsCollection.init(addon, &am.mongoStore)

	am.fileStore.RootPath = path.Join(config.dataPath, "addons")
	am.fileStore.TrashPath = path.Join(config.dataPath, "trash", "addons")
	return
}
func (am *AddonsMan) query(request *Request) (response *QueryResponse) { return }
func (am *AddonsMan) get(request *Request) (response *QueryResponse)   { return }
func (am *AddonsMan) create(request *Request) (response Response)      { return }
func (am *AddonsMan) update(request *Request) (response Response)      { return }
func (am *AddonsMan) delete(request *Request) (response Response)      { return }
