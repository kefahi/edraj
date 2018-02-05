package main

import "path"

// WorkgroupMan to manage the workgroups
type WorkgroupMan struct {
	mongoStore       MongoStore
	addonsCollection MongoCollection

	fileStore Storage
}

func (wgm *WorkgroupMan) init(config *Config) (err error) {
	wgm.mongoStore.init(config.mongoAddress, addon)

	wgm.addonsCollection = MongoCollection{}
	wgm.addonsCollection.init(addon, &wgm.mongoStore)

	wgm.fileStore.RootPath = path.Join(config.dataPath, addon)
	wgm.fileStore.TrashPath = path.Join(config.dataPath, "trash", addon)
	return
}
func (wgm *WorkgroupMan) query(request *Request) (response *QueryResponse) { return }
func (wgm *WorkgroupMan) get(request *Request) (response *QueryResponse)   { return }
func (wgm *WorkgroupMan) create(request *Request) (response Response)      { return }
func (wgm *WorkgroupMan) update(request *Request) (response Response)      { return }
func (wgm *WorkgroupMan) delete(request *Request) (response Response)      { return }
