package main

// MinerMan manage the local miner (this is more of an index for the local content store that keeps exploring the data for more)
type MinerMan struct {
	mongoStore      MongoStore
	minerCollection MongoCollection
}

/*

func (mm *MinerMan) init(mongoAddress string) {
	mm.mongoStore.init(mongoAddress, "miner")
	mm.minerCollection.init("miner", &mm.mongoStore)
}*/

func (mm *MinerMan) init(config *Config) (err error)                  { return }
func (mm *MinerMan) query(request *Request) (response *QueryResponse) { return }
func (mm *MinerMan) get(request *Request) (response *QueryResponse)   { return }
func (mm *MinerMan) create(request *Request) (response Response)      { return }
func (mm *MinerMan) update(request *Request) (response Response)      { return }
func (mm *MinerMan) delete(request *Request) (response Response)      { return }
