package main

// Each storage object maps to a mongodb collection (so its flattened):
// mainly hosting the meta data, as payload remains on the file system

// Identity minimal user/agent pointer
type Identity struct {
	ID          string `bson:"_id" json:"id"`
	Displayname string
	Shortname   string
	Type        string            // Actor, Workgroup, Domain
	PublicKeys  map[string]string // Unique names of the list of active /
	// verifiable public keys
	Domain string
}

// Signature a digital signature of data
type Signature struct {
	Identity     Identity // The actor who signed
	Signature    string
	Keyname      string // The public key used
	FieldsSigned []string
}

// GeoPoint long/lat
type GeoPoint struct {
	Latitude  string
	Longitude string
}

// Attachment file meta
type Attachment struct {
	ID         string `bson:"_id" json:"id"`
	Name       string
	Checksum   string
	Size       int64
	Signature  *Signature `bson:",omitempty" json:",omitempty"` // The Author's signature.
	MediaType  string     // MediaTypes
	References []string   `bson:",omitempty" json:",omitempty"` // TODO fixme list of id's of Contents/Messages
	// that point to this attachment. auto-garbage collection?

	// Ome of the following.
	Path   string `bson:",omitempty" json:",omitempty"` // Local filesystem path
	Binary string `bson:",omitempty" json:",omitempty"` // Binary data
	Text   string `bson:",omitempty" json:",omitempty"` // Text/Markdown/Json
	URL    string `bson:",omitempty" json:",omitempty"` // Remote URL
}

// ActorGroup a logical grouping of actors for esier management
type ActorGroup struct {
	ID      string `bson:"_id" json:"id"`
	Name    string
	Members []Identity `bson:",omitempty" json:",omitempty"`
}

// Permission a permission for a user or a group
type Permission struct {
	Type   string       // admin, read, write, delete, query, manage_reactions
	Actors []Actor      `bson:",omitempty" json:",omitempty"`
	Groups []ActorGroup `bson:",omitempty" json:",omitempty"`
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
	ID             string `bson:"_id" json:"id"`
	Displayname    string
	Shortname      string // unique
	Domain         string //
	Keys           []Keypair
	Address        *Address     `bson:",omitempty" json:",omitempty"`
	Organizations  []string     `bson:",omitempty" json:",omitempty"` // that the user relates to, like work, ngo ...
	Comms          []Contact    `bson:",omitempty" json:",omitempty"` // The user's communication channels
	Biography      string       `bson:",omitempty" json:",omitempty"` // The user's biography
	Contacts       []Actor      `bson:",omitempty" json:",omitempty"` // Firends / Contacts and followed actors
	BannedContacts []Actor      `bson:",omitempty" json:",omitempty"` // Withwhome the user doesn't want to accept or interact
	ActorGroups    []ActorGroup `bson:",omitempty" json:",omitempty"` // Grouping of contacts from the Contacts array
}

// Keypair : PKI keypair
type Keypair struct {
	Public  string
	Private string
}

// Domain a logical pool of users. Future: A domain could have
// multiple legs on more than one server. more of a replica setup.
type Domain struct {
	ID          string `bson:"_id" json:"id"`
	Displayname string
	Shortname   string
	Keys        []Keypair
	ActiveIPs   []string
}

// Page a collection of layed out blocks
type Page struct {
	ID string `bson:"_id" json:"id"`
}

// Site a layout organization of the site and the specific pages it would render
type Site struct {
	ID string `bson:"_id" json:"id"`
} // Site Layout

// Block in a page
type Block struct {
	ID string `bson:"_id" json:"id"`
}

// Addon aka module /plugin
type Addon struct {
	ID   string `bson:"_id" json:"id"`
	Name string
}

// Scheme holdes the details of a schema definition
type Scheme struct {
	ID string `bson:"_id" json:"id"`
}

// Notification of an event
type Notification struct {
	ID        string `bson:"_id" json:"id"`
	Timestamp string
}

// Message aka instant-message / email
type Message struct {
	ID          string `bson:"_id" json:"id"`
	ThreadID    string `bson:"thread_id" json:"thread_id"`
	Recipients  []Actor
	Sender      Actor
	Signature   Signature
	Timestamp   string
	Subject     string
	Body        string
	Attachments []Attachment `bson:",omitempty" json:",omitempty"`
}

// Workgroup Topic-centric as opposed to user-centric
// I.e. its a sepcial type of information repo.
// Think Forums, News websites, ...
type Workgroup struct {
	ID          string `bson:"_id" json:"id"`
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
	ID          string `bson:"_id" json:"id"`
	Shortname   string
	Displayname string
	Timestamp   string
	Signature   Signature
	Geo         *GeoPoint `bson:",omitempty" json:",omitempty"`
	Title       string
	Tags        []string     `bson:",omitempty" json:",omitempty"`
	Categories  []string     `bson:",omitempty" json:",omitempty"`
	Permissions []Permission `bson:",omitempty" json:",omitempty"`
	History     []Change     `bson:",omitempty" json:",omitempty"`
	Body        string
	Attachments []Attachment `bson:",omitempty" json:",omitempty"`
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
	ID              string `bson:"_id" json:"id"`
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
	ID        string `bson:"_id" json:"id"`
	OwnerID   Identity
	Signature Signature
	Timestamp string
	Geo       GeoPoint
	Type      string // like/dislike/...
}

// Signature of data
/*type Signature struct {
	ActorID          string
	ActorDisplayname string
	ActorShortname   string
	ActorType        string // Actor, Workgroup, Domain
	ActorDomain      string
	Signature        string
	PublickeyUsed    string
	FieldsSigned     []string
}*/

// EntryQuery the query object.
type EntryQuery struct {
	EntryType  string `bson:",omitempty" json:",omitempty"` // Of EntryTypes
	Text       string `bson:",omitempty" json:",omitempty"` // free text search
	Date       string `bson:",omitempty" json:",omitempty"` // from-, -to, from-to
	Sort       string `bson:",omitempty" json:",omitempty"` // Sort by fields
	Owner      string `bson:",omitempty" json:",omitempty"` // by ownerid
	Tags       string `bson:",omitempty" json:",omitempty"` // T1,+T2,-T3
	Categories string `bson:",omitempty" json:",omitempty"` // C1,+C2,-C3
	Fields     string `bson:",omitempty" json:",omitempty"` // A,+B,-C
	Offset     int    `bson:",omitempty" json:",omitempty"` //
	Limit      int    `bson:",omitempty" json:",omitempty"` // aka page-size
}

// Entry general entry data
type Entry struct {
	//ID string
	ID string `bson:"_id" json:"id"`

	// Author/owner's identity and proof: signatory
	Signature *Signature `bson:",omitempty" json:",omitempty"` // Author / owener/creator signature
	Timestamp string
	Further   []struct{} `bson:",omitempty" json:",omitempty"` // Further entries to explore.
	// Children/related/trending/top/popular

	Type string // from EntryTypes
	// json with type-specific fields
	Reaction  *Reaction  `bson:",omitempty" json:",omitempty"`
	Comment   *Comment   `bson:",omitempty" json:",omitempty"`
	Content   *Content   `bson:",omitempty" json:",omitempty"`
	Container *Container `bson:",omitempty" json:",omitempty"`
	Message   *Message   `bson:",omitempty" json:",omitempty"`
	Scheme    *Scheme    `bson:",omitempty" json:",omitempty"`
	Workgroup *Workgroup `bson:",omitempty" json:",omitempty"`
	Page      *Page      `bson:",omitempty" json:",omitempty"`
	Block     *Block     `bson:",omitempty" json:",omitempty"`
	Addon     *Addon     `bson:",omitempty" json:",omitempty"`
	Actor     *Actor     `bson:",omitempty" json:",omitempty"`
	Domain    *Domain    `bson:",omitempty" json:",omitempty"`
	// ...
}

// Request object
type Request struct {
	// The Envelope (Requestor details)
	// The subject
	Signature Signature // Requestor's signature
	Timestamp string

	// Action/verb/affordance
	Verb string // query, get,update, create, delete

	//Type       string //Payload type: id string, Entry, EntryQuery
	ObjectType string // EntryTypes

	Entry      *Entry      `bson:",omitempty" json:",omitempty"` // for create
	EntryID    string      `bson:",omitempty" json:",omitempty"` // for get, update, delete
	EntryQuery *EntryQuery `bson:",omitempty" json:",omitempty"` // For query
}

// Response of an api call
type Response struct {
	Status string // succeeded / failed
	Code   int    // Http: 200 OK, 202 Created, 404 Not found,
	// 500 internal server error
	Message string // in case failed the error message is provided
}

// QueryResponse of the Entry api
type QueryResponse struct {
	Status string // succeeded / failed
	Code   int    // Http: 200 OK, 202 Created, 404 Not found,
	// 500 internal server error
	Message  string // in case failed the error message is provided
	Total    int64
	Returned int64
	Entries  []Entry `json:"entries,omitempty"`
}
