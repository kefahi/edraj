package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	actor        = "actor"
	addon        = "addon"
	attachment   = "attachment"
	block        = "block"
	comment      = "comment"
	container    = "container"
	content      = "content"
	crawler      = "crawler"
	domain       = "domain"
	message      = "message"
	miner        = "miner"
	notification = "notification"
	page         = "page"
	reaction     = "reaction"
	schema       = "schema"
	workgroup    = "workgroup"

	trash = "trash"

	succeeded = "succeeded"
	failed    = "failed"

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
		actor:        true,
		addon:        true,
		attachment:   true,
		block:        true,
		comment:      true,
		container:    true,
		content:      true,
		crawler:      true,
		domain:       true,
		message:      true,
		miner:        true,
		notification: true,
		page:         true,
		reaction:     true,
		schema:       true,
		workgroup:    true,
	}
)

// EntryService serving
type EntryService struct {
	//managers map[string]Manager
	manager Manager
}

// Manager interface
type Manager interface {
	init(config *Config) (err error)
	query(request *Request) (response *Response)
	get(request *Request) (response *Response)
	create(request *Request) (response *Response)
	update(request *Request) (response *Response)
	delete(request *Request) (response *Response)
}

func (es *EntryService) init(config Config) (err error) {

	es.manager = &DefaultMan{}
	es.manager.init(&config)
	/*
		es.managers = map[string]Manager{}
		es.managers[actor] = &defaultMan
		es.managers[workgroup] = &defaultMan
		es.managers[domain] = &defaultMan
		es.managers[addon] = &defaultMan
		es.managers[content] = &defaultMan
		es.managers[message] = &defaultMan
		es.managers[schema] = &defaultMan
		es.managers[crawler] = &defaultMan
		es.managers[notification] = &defaultMan
	*/
	/*
		es.managers[workgroup] = &WorkgroupMan{}
		err = es.managers[workgroup].init(&config)
		es.managers[actor] = &ActorMan{}
		err = es.managers[actor].init(&config)
		es.managers[addon] = &DefaultMan{}
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
	*/
	return
}

func respond(w http.ResponseWriter, response *Response) {
	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	w.Write(data)
}

// EntryAPI serves the api
// TODO validate requester id (by checking the signature)
// TODO validate access-control
func EntryAPI(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request Request
	body, _ := ioutil.ReadAll(r.Body)
	//if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	if err := json.Unmarshal(body, &request); err != nil {
		respond(w, &Response{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	//var data bytes.Buffer
	//var me map[string]interface{}
	//body2, _ := json.Marshal(request)
	//json.Unmarshal(body2, &me)
	//printme(me, "", &data)

	//fmt.Println(data.String())

	// TODO validate request
	/*manager, ok := entryService.managers[request.ObjectType]
	if !ok {
		respond(w, Response{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Invalid entry type requested: " + request.ObjectType,
		})
		return
	}*/

	switch strings.ToUpper(request.Verb) {
	case query:
		respond(w, entryService.manager.query(&request))
	case get:
		respond(w, entryService.manager.get(&request))
	case create:
		respond(w, entryService.manager.create(&request))
	case update:
		respond(w, entryService.manager.update(&request))
	case delete:
		respond(w, entryService.manager.delete(&request))
	}
}

// NotImplementedAPI ...
func NotImplementedAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not implemented yet")
}

// HelloAPI ...
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

/*
// Hello (EHLO message)
func (server *DomainsMan) Hello(iIn Identity, sIn Signature, messageIn string) (iOut Identity,
	sOut Signature, messageOut string, token string) {
	// ? Is there a need to establish some form of a basic session here, after verification of the signature?
	return Identity{}, Signature{}, "Hello back how can I help you?",
		"token of verification - alternative to session"
}*/

/*
	Verbs: : PUT (Create), GET (Query ?query=&fields=), POST (Update), DELETE (Delete)
	All ids / shortnames are unique across the domain: content, container, identity,
	Information System: Content/Container self-described

	  QueryAPI Returns links to a sub-set of subentries (with pagination)
In-memory or ondisk collection of ids vs type/location so entries can be easily retrieved.
Query: q=type:[domain,actor,message,content,container,attachment,comment],
         text:,date:,owner:,id:,tags:,categories: (Url param or get-in-body-request)

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
