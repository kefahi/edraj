package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	domain       = "domain"
	actor        = "actor"
	message      = "message"
	container    = "container"
	content      = "content"
	attachment   = "attachment"
	comment      = "comment"
	schema       = "schema"
	addon        = "addon"
	miner        = "miner"
	crawler      = "crawler"
	notification = "notification"
	workgroup    = "workgroup"
	trash        = "trash"

	// Types of exposed verbs
	query  = "QUERY"
	get    = "GET"
	create = "CREATE"
	update = "UPDATE"
	delete = "DELETE"
)

var (
	// EntryTypes the "enum" of accepted entrytypes
	EntryTypes = map[string]bool{
		workgroup:  true,
		domain:     true,
		actor:      true,
		message:    true,
		container:  true,
		content:    true,
		attachment: true,
		comment:    true,
		schema:     true,
		addon:      true,
		miner:      true,
		crawler:    true,
	}
)

// EntryService serving
type EntryService struct {
	managers map[string]Manager
}

// Manager interface
type Manager interface {
	init(config *Config) (err error)
	query(request *Request) (response *QueryResponse)
	get(request *Request) (response *QueryResponse)
	create(request *Request) (response Response)
	update(request *Request) (response Response)
	delete(request *Request) (response Response)
}

func (es *EntryService) init(config Config) (err error) {
	es.managers = map[string]Manager{}
	es.managers[workgroup] = &WorkgroupMan{}
	err = es.managers[workgroup].init(&config)
	es.managers[actor] = &ActorMan{}
	err = es.managers[actor].init(&config)
	es.managers[addon] = &AddonsMan{}
	err = es.managers[addon].init(&config)
	es.managers[domain] = &DomainsMan{}
	err = es.managers[domain].init(&config)
	es.managers[content] = &ContentMan{}
	err = es.managers[content].init(&config)
	es.managers[message] = &MessagesMan{}
	err = es.managers[message].init(&config)
	es.managers[schema] = &SchemaMan{}
	err = es.managers[schema].init(&config)
	es.managers[miner] = &MinerMan{}
	err = es.managers[miner].init(&config)
	es.managers[crawler] = &CrawlersMan{}
	err = es.managers[crawler].init(&config)
	es.managers[notification] = &NotificationsMan{}
	err = es.managers[notification].init(&config)

	return
}

func respond(w http.ResponseWriter, data interface{}) {
	var code int
	switch cast := data.(type) {
	case *Response:
		code = cast.Code
	case *QueryResponse:
		code = cast.Code
	}
	response, _ := json.Marshal(data)
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
		respond(w, response)
		return
	}

	// TODO validate request
	manager, ok := entryService.managers[request.Entry.EntryType]
	if !ok {
		respond(w, Response{
			Status:       "failed",
			Code:         http.StatusBadRequest,
			ErrorMessage: "Invalid entry type requested: " + request.EntryType,
		})
		return
	}

	switch strings.ToUpper(request.RequestType) {
	case query:
		respond(w, manager.query(&request))
	case get:
		respond(w, manager.get(&request))
	case create:
		respond(w, manager.create(&request))
	case update:
		respond(w, manager.update(&request))
	case delete:
		respond(w, manager.delete(&request))
	}
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
