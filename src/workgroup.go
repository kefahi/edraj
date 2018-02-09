package main

import (
	"path"

	mgo "gopkg.in/mgo.v2"
)

// WorkgroupMan to manage the workgroups
type WorkgroupMan struct {
	mongoDb *mgo.Database

	fileStore Storage
}

func (wgm *WorkgroupMan) init(config *Config) (err error) {
	wgm.mongoDb = mongoSession.DB(workgroup)
	wgm.fileStore.RootPath = path.Join(config.dataPath, workgroup)
	wgm.fileStore.TrashPath = path.Join(config.dataPath, trash, workgroup)
	return
}
func (wgm *WorkgroupMan) query(request *Request) (response *QueryResponse) { return }
func (wgm *WorkgroupMan) get(request *Request) (response *QueryResponse)   { return }
func (wgm *WorkgroupMan) create(request *Request) (response Response)      { return }
func (wgm *WorkgroupMan) update(request *Request) (response Response)      { return }
func (wgm *WorkgroupMan) delete(request *Request) (response Response)      { return }
