package main

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	trash = "trash"
	edraj = "edraj"
)

// EntryMan ...
type EntryMan struct {
	mongoSession *mgo.Session
	mongoDb      *mgo.Database
	fileStore    Storage
}

func (man *EntryMan) init(config *Config) (err error) {
	man.mongoSession, err = mgo.Dial(config.mongoAddress)
	if err != nil {
		return
	}
	man.mongoDb = man.mongoSession.DB(edraj)
	man.fileStore.RootPath = path.Join(config.dataPath, edraj)
	man.fileStore.TrashPath = path.Join(config.dataPath, trash, edraj)
	return
}

func entryObject(objectType string, e *Entry, createIfNil bool) (doc interface{}) {
	fieldName := strings.Title(objectType)
	field := reflect.ValueOf(e).Elem().FieldByName(fieldName)
	doc = field.Interface()
	if createIfNil && field.IsNil() {
		field.Set(reflect.Indirect(reflect.New(field.Type().Elem())).Addr())
		doc = field.Interface()
	}
	return
}

func (man *EntryMan) create(request *EntryRequest) (response *Receipt, err error) {
	response = &Receipt{Status: &Status{}}
	//glog.Info(request)
	if request.Entry == nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}
	entryType := strings.ToLower(request.Entry.Type.String())

	doc := entryObject(entryType, request.Entry, false)

	err = man.mongoDb.C(entryType).Insert(doc)
	if err != nil {
		response.Status.Code = int32(codes.Internal)
		response.Status.Message = fmt.Sprintf("Failed to create entry (%v).", request.Entry)
		err = status.Errorf(codes.Internal, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

func (man *EntryMan) update(request *EntryRequest) (response *Receipt, err error) {
	response = &Receipt{Status: &Status{}}
	//glog.Info(request)
	if request.Entry == nil || request.Entry.Id == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.Internal, "Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.Entry.Type.String())

	doc := entryObject(entryType, request.Entry, false)

	err = man.mongoDb.C(entryType).Update(request.Entry.Id, doc)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to update entry (%v).", request.Entry)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

func (man *EntryMan) query(request *QueryRequest) (response *Response, err error) {
	response = &Response{Status: &Status{}}
	//glog.Info(request)

	if request.Query == nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "Query details are missing (%v).", request.Query)
		response.Status.Message = err.Error()
		return
	}
	// TODO consume the Query data (filters, pagination ...etc)
	query := bson.M{}

	response.Entries = []*Entry{}
	entryType := strings.ToLower(request.Query.EntryType.String())

	fieldName := strings.Title(entryType)
	fieldType := reflect.New(reflect.TypeOf(Entry{})).Elem().FieldByName(fieldName).Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(fieldType), 0, 0)
	objects := reflect.New(slice.Type())
	objects.Elem().Set(slice)
	if request.Query.Limit == 0 {
		request.Query.Limit = 2
	}
	err = man.mongoDb.C(entryType).Find(query).Skip(int(request.Query.Offset)).Limit(int(request.Query.Limit)).All(objects.Interface())

	for i := 0; i < objects.Elem().Len(); i++ {
		entry := &Entry{}
		reflect.ValueOf(entry).Elem().FieldByName(fieldName).Set(reflect.ValueOf(objects.Elem().Index(i).Addr().Interface()))
		response.Entries = append(response.Entries, entry)
	}

	if err != nil {
		response.Status.Code = int32(codes.Internal)
		err = status.Errorf(codes.Internal, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

func (man *EntryMan) get(request *IdRequest) (response *Response, err error) {
	//glog.Info(request)
	response = &Response{Status: &Status{}}

	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.EntryType.String())
	entry := Entry{}
	doc := entryObject(entryType, &entry, true)

	err = man.mongoDb.C(entryType).FindId(request.EntryId /*bson.ObjectIdHex(request.EntryID)*/).One(doc)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to get entry (%v).", request.EntryId)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	response.Entries = []*Entry{&entry}
	response.Returned = 1
	response.Total = 1

	return
}

func (man *EntryMan) delete(request *IdRequest) (response *Receipt, err error) {
	// glog.Info(request)
	response = &Receipt{Status: &Status{}}
	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = status.Errorf(codes.InvalidArgument, "EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.EntryType.String())
	err = man.mongoDb.C(entryType).Remove(&struct{ _id string }{_id: request.EntryId})
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to delete entry (%v).", request.EntryId)
		err = status.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}
