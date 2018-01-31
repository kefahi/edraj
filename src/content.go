package main

import "path"

// ContentMan : Content Manager
type ContentMan struct {
	mongoStore            MongoStore
	contentCollection     MongoCollection
	containersCollection  MongoCollection
	attachmentsCollection MongoCollection
	fileStore             Storage
}

func (cm *ContentMan) init(mongoAddress string, rootPath string) {
	cm.mongoStore.init(mongoAddress, "content")
	cm.fileStore.RootPath = path.Join(rootPath, "content")
	cm.fileStore.TrashPath = path.Join(rootPath, "trash", "content")
	cm.contentCollection.init("content", &cm.mongoStore)
}

// NewContainer : Creates a new Container (aka folder)
func (cm *ContentMan) NewContainer(container Content) {}

// Move : Moves Content/Containers around
func (cm *ContentMan) Move(id string, to string) {}

// GetRootContainer returns the the root container
func (cm *ContentMan) GetRootContainer() Container {
	return Container{}
}

// List child-ids by parent
func (cr *Container) List(parentID string) ([]string, error) { return []string{}, nil }

// Delete : deletes a content/container by their uuid (moves to trash)
func (cm *ContentMan) Delete(id string) {}

// Update : updates details
func (c *Content) Update(fields map[string]string) {}

// GetAttachment retrieve the payload
func (c *Content) GetAttachment(attachmentID string) {}

// PutAttachment retrieve the payload
func (c *Content) PutAttachment(contentID, attachmentID string, attachment string) {}

// Update : updates details
func (cr *Container) Update(fields map[string]string) {}

// NewContent : Creates a new Content
func (cr *Container) NewContent(content Content) {}

// Content / Container:
// UpdateMeta/Put/Get (Query is left for the Miner)
// Set permission/Tags/Categories/Description/Notes
