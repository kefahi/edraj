package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	edraj = "edraj"
)

// DefaultMan to manage the edraj installed on the system
type DefaultMan struct {
	mongoDb   *mgo.Database
	fileStore Storage
}

func (man *DefaultMan) init(config *Config) (err error) {
	man.mongoDb = mongoSession.DB(edraj)
	man.fileStore.RootPath = path.Join(config.dataPath, edraj)
	man.fileStore.TrashPath = path.Join(config.dataPath, trash, edraj)
	return
}

func (man *DefaultMan) query(request *Request) (response *QueryResponse) {
	response = &QueryResponse{}
	if !EntryTypes[request.EntryType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request objecttype must have valid values  (%s) (%s)", request.EntryID, request.EntryType)
		return
	}

	response.Entries = []Entry{}

	var err error
	switch request.EntryType {
	case content:
		objects := []Content{}
		err = man.mongoDb.C(request.EntryType).Find(bson.M{}).All(&objects)
		for _, object := range objects {
			value := object
			response.Entries = append(response.Entries, Entry{Content: &value})
		}
	case container:
		objects := []Container{}
		err = man.mongoDb.C(request.EntryType).Find(bson.M{}).All(&objects)
		for _, object := range objects {
			value := object
			response.Entries = append(response.Entries, Entry{Container: &value})
		}

	}

	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to run query (%v). err: %s", request, err.Error())
		return
	}

	response.Returned = int64(len(response.Entries))
	response.Status = succeeded
	response.Code = http.StatusFound

	return
}

func (man *DefaultMan) get(request *Request) (response *QueryResponse) {
	response = &QueryResponse{}

	if request.EntryID == "" || !EntryTypes[request.EntryType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request entryid and/or objecttype must have valid values  (%s) (%s)", request.EntryID, request.EntryType)
		return
	}
	var err error
	entry := Entry{}
	doc := entryObject(request.EntryType, &entry, true)

	err = man.mongoDb.C(request.EntryType).FindId(request.EntryID /*bson.ObjectIdHex(request.EntryID)*/).One(doc)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to retrive entry (%s). err: %s", request.EntryID, err.Error())
		return
	}

	response.Entries = []Entry{entry}
	response.Status = succeeded
	response.Code = http.StatusFound
	response.Returned = 1
	response.Total = 1

	return
}

func entryObject(objectType string, e *Entry, createIfNil bool) (doc interface{}) {
	switch objectType {
	case actor:
		if createIfNil && e.Actor == nil {
			e.Actor = &Actor{}
		}
		doc = e.Actor
	case addon:
		if createIfNil && e.Addon == nil {
			e.Addon = &Addon{}
		}
		doc = e.Addon
	case block:
		if createIfNil && e.Block == nil {
			e.Block = &Block{}
		}
		doc = e.Block
	case comment:
		if createIfNil && e.Comment == nil {
			e.Comment = &Comment{}
		}
		doc = e.Comment
	case container:
		if createIfNil && e.Container == nil {
			e.Container = &Container{}
		}
		doc = e.Container
	case content:
		if createIfNil && e.Content == nil {
			e.Content = &Content{}
		}
		doc = e.Content
	case domain:
		if createIfNil && e.Domain == nil {
			e.Domain = &Domain{}
		}
		doc = e.Domain
	case message:
		if createIfNil && e.Message == nil {
			e.Message = &Message{}
		}
		doc = e.Message
	case page:
		if createIfNil && e.Page == nil {
			e.Page = &Page{}
		}
		doc = e.Page
	case reaction:
		if createIfNil && e.Reaction == nil {
			e.Reaction = &Reaction{}
		}
		doc = e.Reaction
	case schema:
		if createIfNil && e.Schema == nil {
			e.Schema = &Schema{}
		}
		doc = e.Schema
	case workgroup:
		if createIfNil && e.Workgroup == nil {
			e.Workgroup = &Workgroup{}
		}
		doc = e.Workgroup
	default:
		log.Println("Bad object type ", objectType)
	}
	return
}

func (man *DefaultMan) create(request *Request) (response *Response) {
	response = &Response{}
	if request.Entry == nil || !EntryTypes[request.EntryType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entrytype  must have valid value  (%s)", request.EntryType)
		return
	}

	doc := entryObject(request.EntryType, request.Entry, false)

	err := man.mongoDb.C(request.EntryType).Insert(doc)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to create entry (%v). err: %s", request.Entry, err.Error())
		return
	}

	response.Status = succeeded
	response.Code = http.StatusCreated
	return
}

func (man *DefaultMan) update(request *Request) (response *Response) {
	response = &Response{}
	if request.EntryID == "" || request.Entry == nil || !EntryTypes[request.EntryType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.Entry must have valid value  (%v)", request.Entry)
		return
	}

	doc := entryObject(request.EntryType, request.Entry, false)

	err := man.mongoDb.C(request.EntryType).Update(request.EntryID, doc)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to update entry (%v). err: %s", request.Entry, err.Error())
	}

	response.Status = succeeded
	response.Code = http.StatusCreated
	return
}

func (man *DefaultMan) delete(request *Request) (response *Response) {
	response = &Response{}
	if request.EntryID == "" || !EntryTypes[request.EntryType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entryid must have valid value  (%s)", request.EntryID)
		return
	}

	err := man.mongoDb.C(request.EntryType).Remove(&struct{ _id string }{_id: request.EntryID})
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to retrive entry (%s). err: %s", request.EntryID, err.Error())
		return
	}

	response.Status = succeeded
	response.Code = http.StatusGone

	return
}
