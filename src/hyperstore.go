package main

// This is a mongo-based hyper store backed by storage.go file backend.

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Server : Edraj server setup, can host multiple domain-legs
type Server struct {
	publicIPs  []string
	privateIPs []string
}

// Hyperstore : Mongo + filesystem
type Hyperstore struct {
	MongoAddress   string
	DatabaseName   string
	CollectionName string

	session    *mgo.Session
	database   *mgo.Database
	collection *mgo.Collection
}

// Connect to mongodb
func (m *Hyperstore) Connect() {
	var err error
	m.session, err = mgo.Dial(m.MongoAddress)
	if err != nil {
		log.Fatal(err)
	}
	m.database = m.session.DB(m.DatabaseName)
	m.collection = m.database.C(m.CollectionName)
}

// Query matches
func (m *Hyperstore) Query() (interface{}, error) {
	var result interface{}
	err := m.collection.Find(bson.M{}).All(&result)
	return result, err
}

// GetByID for a single document
func (m *Hyperstore) GetByID(id string) (interface{}, error) {
	var result interface{}
	err := m.collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err

}

// Create object
func (m *Hyperstore) Create(docs ...interface{}) error {
	return m.collection.Insert(docs)
}

// Delete object
func (m *Hyperstore) Delete(selector interface{}) error {
	return m.collection.Remove(selector)
}

// Update object
func (m *Hyperstore) Update(id string, selector interface{}) error {
	return m.collection.Update(id, selector)
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
