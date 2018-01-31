package main

// NotificationsMan notifications
type NotificationsMan struct {
	mongoStore              MongoStore
	notificationsCollection MongoCollection
}

func (nm *NotificationsMan) init(mongoAddress string) {
	nm.mongoStore.init(mongoAddress, "notifications")
	nm.notificationsCollection.init("notifications", &nm.mongoStore)
}
