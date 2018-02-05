package main

// This is a mongo-based hyper store backed by storage.go file backend.

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStore : Mongo + filesystem
type MongoStore struct {
	session  *mgo.Session
	database *mgo.Database
}

// MongoCollection handles one collection in a mongodatabase
type MongoCollection struct {
	collection *mgo.Collection
}

// Connect to mongodb
func (m *MongoStore) init(mongoAddress string, dbname string) {
	var err error
	m.session, err = mgo.Dial(mongoAddress)
	if err != nil {
		log.Fatal(err)
	}
	m.database = m.session.DB(dbname)
}

func (c *MongoCollection) init(name string, mongoStore *MongoStore) {
	c.collection = mongoStore.database.C(name)
}

// Query matches
func (c *MongoCollection) Query() (interface{}, error) {
	var result interface{}
	err := c.collection.Find(bson.M{}).All(&result)
	return result, err
}

// GetByID for a single document
func (c *MongoCollection) GetByID(id string) (interface{}, error) {
	var result interface{}
	err := c.collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err

}

// Create object
func (c *MongoCollection) Create(doc interface{}) error {
	switch t := doc.(type) {
	case *Domain:

	default:
		fmt.Println("pub is of type:", t)
		return fmt.Errorf("Unsupported public key type %T", doc)
	}
	return c.collection.Insert(doc)
}

// Delete object
func (c *MongoCollection) Delete(selector interface{}) error {
	return c.collection.Remove(selector)
}

// Update object
func (c *MongoCollection) Update(id string, selector interface{}) error {
	return c.collection.Update(id, selector)
}

/*
// IdentityMan Actor identity manager
type IdentityMan struct{}

// New : Creates a new Identity (user/actor)
func (im *IdentityMan) New(identity Identity) {}

// Delete : deletes a user by their uuid. Initially delete is much like disable/deactivate.
func (im *IdentityMan) Delete(id string) {}

// Update : updates details
func (im *IdentityMan) Update(identity Identity) {}

// WorkgroupMan manages all workgroups
type WorkgroupMan struct{}

// New : Creates a workgroup
func (wgm *WorkgroupMan) New(workgroup Workgroup) {}

// Delete : deletes a workgroup by their uuid. Initially delete is much like disable/deactivate.
func (wgm *WorkgroupMan) Delete(id string) {}

// Update : updates details
func (wg *Workgroup) Update(workgroup Workgroup) {}

// AddonsMan manages addons
type AddonsMan struct{}

// SchemaMan manages schema
type SchemaMan struct{}

// SiteMan manages site: Layouts, Pages, Block
type SiteMan struct{}

// NotificationMan manages notifications
type NotificationMan struct{}

// MessagingMan manages messaging
type MessagingMan struct{}

// Miner the local miner
type Miner struct{}

// Query Content/Container by various attributes like name, tags, actual content ...

// Crawler the public miner
type Crawler struct{} // aka public miner


import (
	"fmt"
	"log"
)

var (

	// User specific within a domain
	userContent       = Storage{}
	userMessages      = Storage{}
	userNotifications = Storage{}
	userPublic        = Storage{} // Pages, Layouts and Blocks
	userAddons        = Storage{}
	userSchema        = Storage{}
	userDomains       = Storage{} // Caching domain details
	userPublicMiners  = Storage{}
	userLocalMiner    = Storage{}
	userCache         = Storage{}
	userPeople        = Storage{} // Friends, followed

	// Domain-level
	identities = Storage{}
	domains    = Storage{}
	addons     = Storage{}
	miners     = Storage{}
	schema     = Storage{}
	workgroups = Storage{} // Cluster of storages per workgroup
	users      = Storage{} // Cluster of storages per user

)

func models() {
	fmt.Println("hi")
	log.Println("hi")

}*/

/*
func (p *Page) create() {}
func (p *Page) query()  {}
func (p *Page) update() {}
func (p *Page) delete() {}
func (p *Page) render() {}

func (p *Block) create() {}
func (p *Block) query()  {}
func (p *Block) update() {}
func (p *Block) delete() {}
func (p *Block) render() {}

func (p *Site) create() {}
func (p *Site) query()  {}
func (p *Site) update() {}
func (p *Site) delete() {}
func (p *Site) render() {}
*/
