package main

import (
	"path"

	mgo "gopkg.in/mgo.v2"
)

// ActorMan to manage the addons installed on the system
type ActorMan struct {
	mongoDb   *mgo.Database
	fileStore Storage
}

func (am *ActorMan) init(config *Config) (err error) {
	am.mongoDb = mongoSession.DB(actor)
	am.fileStore.RootPath = path.Join(config.dataPath, actor)
	am.fileStore.TrashPath = path.Join(config.dataPath, trash, actor)
	return
}
func (am *ActorMan) query(request *Request) (response *QueryResponse) { return }
func (am *ActorMan) get(request *Request) (response *QueryResponse)   { return }
func (am *ActorMan) create(request *Request) (response Response)      { return }
func (am *ActorMan) update(request *Request) (response Response)      { return }
func (am *ActorMan) delete(request *Request) (response Response)      { return }
