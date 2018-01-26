package main

// Each storage object maps to a mongodb collection (so its flattened):
// mainly hosting the meta data, as payload remains on the file system

import (
	"gopkg.in/mgo.v2/bson"
)

// Identity minimal user/agent pointer
type Identity struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Displayname string
	Shortname   string
	Type        string            // Actor, Workgroup, Domain
	PublicKeys  map[string]string // Unique names of the list of active / verifiable public keys
	Domain      string
}

// Signature a digital signature of data
type Signature struct {
	Identity  Identity // The actor who signed
	Signature string
	Keyname   string
}

// GeoPoint long/lat
type GeoPoint struct {
	Latitude  string
	Longitude string
}

// Attachment file meta
type Attachment struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Name       string
	Checksum   string
	Size       int64
	Signature  Signature // The Author's signature.
	MediaType  string    // MediaTypes
	References []string  // TODO fixme list of id's of Contents/Messages that point to this attachment. auto-garbage collection?

	// Ome of the following.
	Path   string // Local filesystem path
	Binary string // Binary data
	Text   string // Text/Markdown/Json
	URL    string // Remote URL
}

// ActorGroup a logical grouping of actors for esier management
type ActorGroup struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	name    string
	members []Identity
}

// Permission a permission for a user or a group
type Permission struct {
	Type   string // admin, read, write, delete, query, manage_reactions
	Actors []Actor
	Groups []ActorGroup
}

// Contact means of acessing the user or its related info
// Emails
type Contact struct {
	Type  string // Email, Phone, Url, Twitter, Skype, FB, Linked in
	Value string
}

// Address a person's address / geographical locatiotion
type Address struct {
	Geo     GeoPoint
	Line1   string
	Line2   string
	Zipcode string
	City    string
	State   string
	Country string
}

// Actor aka User/Agent
type Actor struct {
	ID               bson.ObjectId `bson:"_id" json:"id"`
	Displayname      string
	Shortname        string // unique
	Domain           string
	Keys             []Keypair
	Address          Address
	Organizations    []string  // that the user relates to, like work, ngo ...
	Comms            []Contact // The user's communication channels
	Biography        string    // The user's biography
	ActorConnections []Actor   // Firends / Contacts and followed actors
	ActorGroups      []ActorGroup
}

// Keypair : PKI keypair
type Keypair struct {
	Public  string
	Private string
}

// Domain a logical pool of users
type Domain struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Displayname string
	Shortname   string
	Keys        []Keypair
	ActiveIPs   []string
}

// Page a collection of layed out blocks
type Page struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
}

// Site a layout organization of the site and the specific pages it would render
type Site struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
} // Site Layout

// Block in a page
type Block struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
}

// Addon aka module /plugin
type Addon struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
}

// Scheme holdes the details of a schema definition
type Scheme struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
}

// Notification of an event
type Notification struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Timestamp string
}

// Message aka instant-message / email
type Message struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	ThreadID    string        `bson:"thread_id" json:"thread_id"`
	Recipients  []Actor
	Sender      Actor
	Signature   Signature
	Timestamp   string
	Subject     string
	Body        string
	Attachments []Attachment
}

// Workgroup Topic-centric as opposed to user-centric
// I.e. its a sepcial type of information repo.
// Think Forums, News websites, ...
type Workgroup struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Shortname   string
	Displayname string
	Domain      string
	Keys        []Keypair
	members     []Identity
	Permissions []Permission
}

// Change change
type Change struct {
	Timestamp string   // When
	Actor     Identity // Who
	Delta     string   // What
}

// Content content
type Content struct {
	ID          string
	Shortname   string
	Displayname string
	Timestamp   string
	Signature   Signature
	Geo         GeoPoint
	Title       string
	Tags        []string
	Categories  []string
	Permissions []Permission
	history     []Change
	Body        string
	Attachments []Attachment
}

// Container of content
type Container struct {
	ID          string
	OwnerID     string
	Shortname   string
	Displayname string
	Timestamp   string
	Tags        []string
	Categories  []string
	Permissions []Permission
}

// Comment ...
type Comment struct {
	ID              bson.ObjectId `bson:"_id" json:"id"`
	Actor           Identity
	Timestamp       string
	Geo             GeoPoint
	Signature       Signature
	ContentID       string
	ParentCommentID string
	Text            string
}

// Reaction / comment should be part of reaction?
type Reaction struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	OwnerID   Identity
	Signature Signature
	Timestamp string
	Geo       GeoPoint
	Type      string // like/dislike/...
}
