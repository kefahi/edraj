package main

import (
	"fmt"
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
	if !EntryTypes[request.ObjectType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request objecttype must have valid values  (%s) (%s)", request.EntryID, request.ObjectType)
		return
	}

	response.Entries = []Entry{}

	var err error
	switch request.ObjectType {
	case content:
		objects := []Content{}
		err = man.mongoDb.C(request.ObjectType).Find(bson.M{}).All(&objects)
		for _, object := range objects {
			value := object
			response.Entries = append(response.Entries, Entry{Content: &value})
		}
	case container:
		objects := []Container{}
		err = man.mongoDb.C(request.ObjectType).Find(bson.M{}).All(&objects)
		for _, object := range objects {
			response.Entries = append(response.Entries, Entry{Container: &object})
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

	if request.EntryID == "" || !EntryTypes[request.ObjectType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request entryid and/or objecttype must have valid values  (%s) (%s)", request.EntryID, request.ObjectType)
		return
	}
	var err error
	entry := Entry{}
	var object interface{}
	switch request.ObjectType {
	case content:
		object = &entry.Content
	case container:
		object = &entry.Container

	}
	err = man.mongoDb.C(request.ObjectType).FindId(request.EntryID /*bson.ObjectIdHex(request.EntryID)*/).One(object)
	response.Entries = []Entry{entry}
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to retrive entry (%s). err: %s", request.EntryID, err.Error())
		return
	}

	response.Status = succeeded
	response.Code = http.StatusFound
	//response.Entries = append(response.Entries, entry)
	response.Returned = 1
	response.Total = 1

	return
}

func (man *DefaultMan) create(request *Request) (response Response) {
	if !EntryTypes[request.ObjectType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entry  must have valid value  (%v)", request.Entry)
		return
	}

	var doc interface{}
	switch request.ObjectType {
	case content:
		doc = request.Entry.Content
	case container:
		doc = request.Entry.Container
	case comment:
		doc = request.Entry.Comment
	case reaction:
		doc = request.Entry.Reaction
	case message:
		doc = request.Entry.Message
	case actor:
		doc = request.Entry.Actor
	case workgroup:
		doc = request.Entry.Workgroup
	case schema:
		doc = request.Entry.Scheme
	}
	err := man.mongoDb.C(request.ObjectType).Insert(doc)
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

func (man *DefaultMan) update(request *Request) (response Response) {
	if request.EntryID == "" || !EntryTypes[request.ObjectType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.Entry must have valid value  (%v)", request.Entry)
		return
	}
	err := man.mongoDb.C(request.ObjectType).Update(request.Entry.ID, &request.Entry)
	if err != nil {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("Failed to update entry (%v). err: %s", request.Entry, err.Error())
	}

	response.Status = succeeded
	response.Code = http.StatusCreated
	return
}

func (man *DefaultMan) delete(request *Request) (response Response) {
	if request.EntryID == "" || !EntryTypes[request.ObjectType] {
		response.Status = failed
		response.Code = http.StatusBadRequest
		response.Message = fmt.Sprintf("request.entryid must have valid value  (%s)", request.EntryID)
		return
	}

	err := man.mongoDb.C(request.ObjectType).Remove(&struct{ _id string }{_id: request.EntryID})
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
