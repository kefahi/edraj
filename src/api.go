package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var (
	// EntryTypes the "enum" of accepted entrytypes
	EntryTypes = map[string]bool{
		"domain":     true,
		"actor":      true,
		"message":    true,
		"container":  true,
		"content":    true,
		"attachment": true,
		"comment":    true,
		"schema":     true,
		"addon":      true,
		"miner":      true,
	}
)

// EntryService serving
type EntryService struct{}

func init() {
	// initialize entry
}

// EntryQuery the query object.
type EntryQuery struct {
	entryType  string // Of EntryTypes
	text       string // free text search
	date       string // from-, -to, from-to
	sort       string // Sort by fields
	owner      string // by ownerid
	tags       string // T1,+T2,-T3
	categories string // C1,+C2,-C3
	fields     string // A,+B,-C
	Offset     int    //
	Limit      int    // aka page-size
}

// Entry general entry data
type Entry struct {
	ID                 string
	ownerID            string
	owerDomain         string
	signature          string
	signaturePublickey string
	timestamp          string
	further            []struct{} // Further entries to explore. Children/related

	subject string

	entryType    string // EntryTypes
	entryPayload string // json with type-specific fields
	// ...
}

// EntryRequest object
type EntryRequest struct {
	requesterID     string // uuid of the requestor
	requesterDomain string //

	requestSignature   string
	signaturePublickey string
	timestamp          string

	requestType string // query, get,update, create, delete
	entryID     string // for get, update, delete
	entry       Entry  // for create

}

// EntryResponse of the Entry api
type EntryResponse struct {
	total    int64
	returned int64
	entries  []Entry
}

// Query : All fields are optional, when empty it returns the root container
func (a *EntryService) query(q *EntryQuery) {}

// Get : Returns one specific entry object
func (a *EntryService) get(id string)              {}
func (a *EntryService) create(e *Entry)            {}
func (a *EntryService) update(id string, e *Entry) {}
func (a *EntryService) delete(id string)           {}

// EntryAPI serves the api
// TODO validate requester id (by checking the signature)
// TODO validate access-control
func EntryAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("Got requst/method", r.URL.Path, r.Method)
	switch r.Method {
	case "POST":
		UpdateEntryAPI(w, r)
	case "PUT":
		CreateEntryAPI(w, r)
	case "GET":
		re, _ := regexp.Compile("/api/entry/(.*)")
		values := re.FindStringSubmatch(r.URL.Path)
		log.Println(values, len(values))
		if len(values) > 1 && values[1] != "" {
			GetEntryAPI(w, r) //values[1]
		} else {
			QueryAPI(w, r)
		}
	case "DELETE":
		DeleteEntryAPI(w, r)
	default:
		// TODO handle this case
	}

}

/*QueryAPI Returns links to a sub-set of subentries (with pagination)
In-memory or ondisk collection of ids vs type/location so entries can be easily retrieved.
Query: q=type:[domain,actor,message,content,container,attachment,comment],
         text:,date:,owner:,id:,tags:,categories: (Url param or get-in-body-request)
Fields: f=a,+b,-c
Offset: o=10
Limit: l=5
Sort: s=type,date,owner

response:
entries:[{type:, id:, further:{}, actor:, sginature:, timestamp:}]
total: nn
returned: kk
*/
func QueryAPI(w http.ResponseWriter, r *http.Request) { log.Println("In queryapi") }

// GetEntryAPI retrievs one entry object
func GetEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in getentryapi") }

// CreateEntryAPI create
func CreateEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in createentryapi") }

// UpdateEntryAPI update
func UpdateEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in updateentryapi") }

// DeleteEntryAPI delete
func DeleteEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in deleteentryapi") }

// NotImplementedAPI ...
func NotImplementedAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not implemented yet")
}

// HelloAPI ...
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

// Hello (EHLO message)
func (server *Server) Hello(iIn Identity, sIn Signature, messageIn string) (iOut Identity,
	sOut Signature, messageOut string, token string) {
	// ? Is there a need to establish some form of a basic session here, after verification of the signature?
	return Identity{}, Signature{}, "Hello back how can I help you?",
		"token of verification - alternative to session"

}

/*
	Verbs: : PUT (Create), GET (Query ?query=&fields=), POST (Update), DELETE (Delete)

	All ids / shortnames are unique across the domain: content, container, identity,


	Information System: Content/Container self-described

	/api/entry/{ID/shortname} <= actor, message, content, container, attachment, history, comment, react
	// Returns links to a sub-set of subentries (with pagination)
	// In-memory or ondisk collection of ids vs type/location so entries can be easily retrieved.
	Query: q=type:[domain,actor,message,content,container,attachment,comment],text:,date:,owner:,id:
	       (Url param or get-in-body-request)
	Fields: f=a,+b,-c
	Offset: o=10
	Limit: l=5
	Sort: s=type,date,owner

	response:
	entries:[{type:, id:, further:{}, actor:, sginature:, timestamp:}]
	total: nn
	returned: kk


	/api/content/{contentID}
	/api/content/{contentID}
	/api/attachment/{ID}
	/api/comment/{ID}
	/api/container/{containerID}
	/api/schema/{shortname/schemaID}
	/api/actor/{shortname/actorID}
	/api/react/{contentID}  POST
	/api/comment/{contentID}
	/api/message/{messageID}  ?threadID=
	/api/notification


	/api/schema
	/api/addon  Query: List/Getone
	/api/domain
	/api/miner ? miners mainly consume/index and offer a search api

	/    web (regular HTML/non-api)
*/
