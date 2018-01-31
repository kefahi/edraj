package main

// MinerMan manage the local miner (this is more of an index for the local content store that keeps exploring the data for more)
type MinerMan struct {
	mongoStore      MongoStore
	minerCollection MongoCollection
}

func (mm *MinerMan) init(mongoAddress string) {
	mm.mongoStore.init(mongoAddress, "miner")
	mm.minerCollection.init("miner", &mm.mongoStore)
}
