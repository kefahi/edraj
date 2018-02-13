package main

// Each storage object maps to a mongodb collection (so its flattened):
// mainly hosting the meta data, as payload remains on the file system

// Identity minimal user/agent pointer
type Identity struct {
	ID          string `bson:"_id"`
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

// Geopoint long/lat
type Geopoint struct {
	Latitude  float64
	Longitude float64
}

// Attachment file meta
type Attachment struct {
	ID         string `bson:"_id"`
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
	ID      string `bson:"_id"`
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
	Geo     *Geopoint `bson:",omitempty" json:",omitempty"`
	Line1   string
	Line2   string
	Zipcode string
	City    string
	State   string
	Country string
}

// Actor aka User/Agent
type Actor struct {
	ID             string `bson:"_id"`
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
	Name    string
	Public  string
	Private string
}

// Domain a logical pool of users. Future: A domain could have
// multiple legs on more than one server. more of a replica setup.
type Domain struct {
	ID          string `bson:"_id"`
	Displayname string
	Shortname   string
	Keys        []Keypair
	ActiveIPs   []string
}

// Page a collection of layed out blocks
type Page struct {
	ID string `bson:"_id"`
}

// Site a layout organization of the site and the specific pages it would render
type Site struct {
	ID string `bson:"_id"`
} // Site Layout

// Block in a page
type Block struct {
	ID string `bson:"_id"`
}

// Addon aka module /plugin
type Addon struct {
	ID   string `bson:"_id"`
	Name string
}

// Miner aka module /plugin
type Miner struct {
	ID   string `bson:"_id"`
	Name string
}

// Crawler aka module /plugin
type Crawler struct {
	ID   string `bson:"_id"`
	Name string
}

// Schema holdes the details of a schema definition
type Schema struct {
	ID string `bson:"_id"`
}

// Notification of an event
type Notification struct {
	ID        string `bson:"_id"`
	Timestamp string
}

// Message aka instant-message / email
type Message struct {
	ID          string `bson:"_id"`
	ThreadID    string `bson:"thread_id" json:"thread_id"`
	Recipients  []Identity
	Signature   Signature // Sender's signature
	Timestamp   string
	Subject     string
	Body        string
	Attachments []Attachment `bson:",omitempty" json:",omitempty"`
	Geopoint    *Geopoint    `bson:",omitempty" json:",omitempty"`
}

// Workgroup Topic-centric as opposed to user-centric
// I.e. its a sepcial type of information repo.
// Think Forums, News websites, ...
type Workgroup struct {
	ID          string `bson:"_id"`
	Shortname   string
	Displayname string
	Domain      string
	Keys        []Keypair
	Members     []Identity
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
	ID            string `bson:"_id"`
	Shortname     string
	Displayname   string
	Timestamp     string
	Signature     Signature
	Geopoint      *Geopoint `bson:",omitempty" json:",omitempty"`
	Title         string
	Tags          []string     `bson:",omitempty" json:",omitempty"`
	Categories    []string     `bson:",omitempty" json:",omitempty"`
	Permissions   []Permission `bson:",omitempty" json:",omitempty"`
	History       []Change     `bson:",omitempty" json:",omitempty"`
	Body          string
	StucturedBody map[string]interface{} `bson:",omitempty" json:",omitempty"`
	Attachments   []Attachment           `bson:",omitempty" json:",omitempty"`
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
	ID              string `bson:"_id"`
	Actor           Identity
	Timestamp       string
	Geo             *Geopoint `bson:",omitempty" json:",omitempty"`
	Signature       Signature
	ContentID       string
	ParentCommentID string
	Text            string
}

// Reaction / comment should be part of reaction?
type Reaction struct {
	ID        string `bson:"_id"`
	OwnerID   Identity
	Signature Signature
	Timestamp string
	Geo       *Geopoint `bson:",omitempty" json:",omitempty"`
	Type      string    // like/dislike/...
}

// Query the query object.
type Query struct {
	//EntryType  []string `bson:",omitempty" json:",omitempty"` // Of EntryTypes
	Text       string `bson:",omitempty" json:",omitempty"` // free text search
	Date       string `bson:",omitempty" json:",omitempty"` // from-, -to, from-to
	Sort       string `bson:",omitempty" json:",omitempty"` // Sort by fields
	Path       string `bson:",omitempty" json:",omitempty"` // Object path. support patterns?
	Owner      string `bson:",omitempty" json:",omitempty"` // by ownerid
	Tags       string `bson:",omitempty" json:",omitempty"` // T1,+T2,-T3
	Categories string `bson:",omitempty" json:",omitempty"` // C1,+C2,-C3
	Fields     string `bson:",omitempty" json:",omitempty"` // A,+B,-C
	Offset     int    `bson:",omitempty" json:",omitempty"` //
	Limit      int    `bson:",omitempty" json:",omitempty"` // aka page-size
}

// Entry general entry data
type Entry struct {
	Further []struct{} `bson:",omitempty" json:",omitempty"` // Further entries to explore. Children/related/trending/top/popular

	// json with type-specific fields
	Actor        *Actor        `bson:",omitempty" json:",omitempty"`
	Addon        *Addon        `bson:",omitempty" json:",omitempty"`
	Attachment   *Attachment   `bson:",omitempty" json:",omitempty"`
	Block        *Block        `bson:",omitempty" json:",omitempty"`
	Comment      *Comment      `bson:",omitempty" json:",omitempty"`
	Container    *Container    `bson:",omitempty" json:",omitempty"`
	Content      *Content      `bson:",omitempty" json:",omitempty"`
	Crawler      *Crawler      `bson:",omitempty" json:",omitempty"`
	Domain       *Domain       `bson:",omitempty" json:",omitempty"`
	Message      *Message      `bson:",omitempty" json:",omitempty"`
	Miner        *Miner        `bson:",omitempty" json:",omitempty"`
	Notification *Notification `bson:",omitempty" json:",omitempty"`
	Page         *Page         `bson:",omitempty" json:",omitempty"`
	Reaction     *Reaction     `bson:",omitempty" json:",omitempty"`
	Schema       *Schema       `bson:",omitempty" json:",omitempty"`
	Workgroup    *Workgroup    `bson:",omitempty" json:",omitempty"`
	// ...
}

// Request object TODO support bulk operations as well
type Request struct {
	Signature Signature // Requestor's signature
	// Action/verb/affordance
	Verb      string // get, query, update, create, delete
	EntryType string `bson:",omitempty" json:",omitempty"` // from EntryTypes

	// Operation scale
	// 'id': just the EntryID field, 'entry': just the Entry field, 'query': Just the Query field, 'entries': Just the Entries field
	Scale   string `bson:",omitempty" json:",omitempty"`
	EntryID string `bson:",omitempty" json:",omitempty"` // for get, delete. scale = id
	Entry   *Entry `bson:",omitempty" json:",omitempty"` // for single create, update. scale = entry
	Query   *Query `bson:",omitempty" json:",omitempty"` // For query. scale = query
}

/*// Response of an api call
type Response struct {
	Status  string // succeeded / failed
	Code    int    // Http: 200 OK, 202 Created, 404 Not found, 500 internal server error
	Message string // in case failed the error message is provided
}*/

// Response of the Entry api
type Response struct {
	Status   string  // succeeded / failed
	Code     int     // Http: 200 OK, 202 Created, 404 Not found, 500 internal server error
	Message  string  `json:",omitempty"` //  message about the operation
	Total    int64   `json:",omitempty"`
	Returned int64   `json:",omitempty"`
	Entries  []Entry `json:",omitempty"`
}
