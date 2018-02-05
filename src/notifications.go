package main

// NotificationsMan notifications
type NotificationsMan struct {
	mongoStore              MongoStore
	notificationsCollection MongoCollection
}

/*

func (nm *NotificationsMan) init(mongoAddress string) {
	nm.mongoStore.init(mongoAddress, "notifications")
	nm.notificationsCollection.init("notifications", &nm.mongoStore)
}*/

func (nm *NotificationsMan) init(config *Config) (err error)                  { return }
func (nm *NotificationsMan) query(request *Request) (response *QueryResponse) { return }
func (nm *NotificationsMan) get(request *Request) (response *QueryResponse)   { return }
func (nm *NotificationsMan) create(request *Request) (response Response)      { return }
func (nm *NotificationsMan) update(request *Request) (response Response)      { return }
func (nm *NotificationsMan) delete(request *Request) (response Response)      { return }
