package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
		"crawler":    true,
	}
)

// EntryService serving
type EntryService struct {
	domainsMan       DomainsMan // This includes the list of other domains and the local domain's identities?
	contentMan       ContentMan // This is for content, container,attachments, comments
	messagesMan      MessagesMan
	schemaMan        SchemaMan
	addonsMan        AddonsMan
	minerMan         MinerMan
	crawlersMan      CrawlersMan // aka public miners
	notificationsMan NotificationsMan
}

func (es *EntryService) init(mongoAddress string, rootDataPath string) {
	es.addonsMan.init(mongoAddress, rootDataPath)
	es.domainsMan.init(mongoAddress)
	es.contentMan.init(mongoAddress, rootDataPath)
	es.messagesMan.init(mongoAddress, rootDataPath)
	es.schemaMan.init(mongoAddress)
	es.minerMan.init(mongoAddress)
	es.crawlersMan.init(mongoAddress)
	es.notificationsMan.init(mongoAddress)
}

// EntryQuery the query object.
type EntryQuery struct {
	EntryType  string // Of EntryTypes
	Text       string // free text search
	Date       string // from-, -to, from-to
	Sort       string // Sort by fields
	Owner      string // by ownerid
	Tags       string // T1,+T2,-T3
	Categories string // C1,+C2,-C3
	Fields     string // A,+B,-C
	Offset     int    //
	Limit      int    // aka page-size
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

// Entry general entry data
type Entry struct {
	ID string

	// Author/owner's identity and proof: signatory
	Signature Signature
	Timestamp string
	Further   []struct{} `json:"further,omitempty"` // Further entries to explore. Children/related/trending/top/popular

	EntryType    string // from EntryTypes
	EntryPayload string // json with type-specific fields
	// ...
}

// Request object
type Request struct {
	// The Envelope (Requestor details)
	// The subject
	Signature Signature
	Timestamp string

	// Action/verb/affordance
	RequestType string // query, get,update, create, delete

	// Object
	// The Body
	// Based on the requestType one of the following will be provided
	EntryID    string     // for get, update, delete
	Entry      Entry      // for create
	EntryQuery EntryQuery // For query
}

// Response of an api call
type Response struct {
	Status       string // succeeded / failed
	Code         int    // Http: 200 OK, 202 Created, 404 Not found, 500 internal server error
	ErrorMessage string // in case failed the error message is provided
}

// QueryResponse of the Entry api
type QueryResponse struct {
	Status       string // succeeded / failed
	Code         int    // Http: 200 OK, 202 Created, 404 Not found, 500 internal server error
	ErrorMessage string // in case failed the error message is provided
	Total        int64
	Returned     int64
	Entries      []Entry `json:"entries,omitempty"`
}

// Query : when empty it returns the root container
func (es *EntryService) query(r *Request) *QueryResponse { return &QueryResponse{} }

// Get : Returns one specific entry object based on the provided id/shortname
func (es *EntryService) get(r *Request) *QueryResponse { return &QueryResponse{} }
func (es *EntryService) create(r *Request) Response    { return Response{} }
func (es *EntryService) update(r *Request) Response    { return Response{} }
func (es *EntryService) delete(r *Request) Response    { return Response{} }

func respond(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// EntryAPI serves the api
// TODO validate requester id (by checking the signature)
// TODO validate access-control
func EntryAPI(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := Response{
			Status:       "failed",
			Code:         http.StatusBadRequest,
			ErrorMessage: "Invalid request: " + err.Error(),
		}
		respond(w, http.StatusBadRequest, response)
		return
	}

	// TODO validate request

	switch strings.ToUpper(request.RequestType) {
	case "QUERY":
		response := entryService.query(&request)
		respond(w, response.Code, response)
	case "GET":
		response := entryService.get(&request)
		respond(w, response.Code, response)
	case "CREATE":
		response := entryService.create(&request)
		respond(w, response.Code, response)
	case "UPDATE":
		response := entryService.update(&request)
		respond(w, response.Code, response)
	case "DELETE":
		response := entryService.get(&request)
		respond(w, response.Code, response)
	}
	/*
		switch r.Method {
		case "POST":
			UpdateEntryAPI(w, r)
		case "PUT":
			CreateEntryAPI(w, r)
		case "GET":
			re, _ := regexp.Compile("/api/entry/(.*)")
			values := re.FindStringSubmatch(r.URL.Path)
			if len(values) > 1 && values[1] != "" {
				GetEntryAPI(w, r) //values[1]
			} else {
				QueryAPI(w, r)
			}
		case "DELETE":
			DeleteEntryAPI(w, r)
		default:
			// TODO handle this case
		}*/
}

/*  QueryAPI Returns links to a sub-set of subentries (with pagination)
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
/*
func QueryAPI(w http.ResponseWriter, r *http.Request) {
	// TODO contstruct a query object
	log.Println("In queryapi")
}

// GetEntryAPI retrievs one entry object
func GetEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in getentryapi") }

// CreateEntryAPI create
func CreateEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in createentryapi") }

// UpdateEntryAPI update
func UpdateEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in updateentryapi") }

// DeleteEntryAPI delete
func DeleteEntryAPI(w http.ResponseWriter, r *http.Request) { log.Println("in deleteentryapi") }
*/

// NotImplementedAPI ...
func NotImplementedAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not implemented yet")
}

// HelloAPI ...
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

// Hello (EHLO message)
func (server *DomainsMan) Hello(iIn Identity, sIn Signature, messageIn string) (iOut Identity,
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
