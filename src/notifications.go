package main

import mgo "gopkg.in/mgo.v2"

// NotificationsMan notifications
type NotificationsMan struct {
	mongoDb *mgo.Database
}

func (nm *NotificationsMan) init(config *Config) (err error) {
	nm.mongoDb = mongoSession.DB(domain)
	return
}
func (nm *NotificationsMan) query(request *Request) (response *QueryResponse) { return }
func (nm *NotificationsMan) get(request *Request) (response *QueryResponse)   { return }
func (nm *NotificationsMan) create(request *Request) (response Response)      { return }
func (nm *NotificationsMan) update(request *Request) (response Response)      { return }
func (nm *NotificationsMan) delete(request *Request) (response Response)      { return }
