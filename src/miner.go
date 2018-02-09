package main

import mgo "gopkg.in/mgo.v2"

// MinerMan manage the local miner (this is more of an index for the local content store that keeps exploring the data for more)
type MinerMan struct {
	mongoDb *mgo.Database
}

func (mm *MinerMan) init(config *Config) (err error) {
	mm.mongoDb = mongoSession.DB(miner)
	return
}
func (mm *MinerMan) query(request *Request) (response *QueryResponse) { return }
func (mm *MinerMan) get(request *Request) (response *QueryResponse)   { return }
func (mm *MinerMan) create(request *Request) (response Response)      { return }
func (mm *MinerMan) update(request *Request) (response Response)      { return }
func (mm *MinerMan) delete(request *Request) (response Response)      { return }
