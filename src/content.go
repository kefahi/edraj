package main

// ContentMan : Content Manager
type ContentMan struct{}

// NewContainer : Creates a new Container (aka folder)
func (cm *ContentMan) NewContainer(container Content) {}

// Move : Moves Content/Containers around
func (cm *ContentMan) Move(id string, to string) {}

// List child-ids by parent
func (cm *ContentMan) List(parentID string) ([]string, error) { return []string{}, nil }

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
