package main

import "path"

// ActorMan to manage the addons installed on the system
type ActorMan struct {
	mongoStore MongoStore
	fileStore  Storage
}

func (am *ActorMan) init(config *Config) (err error) {
	am.mongoStore.init(config.mongoAddress, addon)
	am.fileStore.RootPath = path.Join(config.dataPath, addon)
	am.fileStore.TrashPath = path.Join(config.dataPath, "trash", addon)
	return
}
func (am *ActorMan) query(request *Request) (response *QueryResponse) { return }
func (am *ActorMan) get(request *Request) (response *QueryResponse)   { return }
func (am *ActorMan) create(request *Request) (response Response)      { return }
func (am *ActorMan) update(request *Request) (response Response)      { return }
func (am *ActorMan) delete(request *Request) (response Response)      { return }
