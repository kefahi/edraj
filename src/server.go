package main

import (
	"flag"

	"fmt"
	"net"
	//"net/http"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	/*
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
	*/
	trash = "trash"

	edraj = "edraj"
)

// EntryServer implements EntryServiceServer
type EntryServer struct {
	mongoSession *mgo.Session
	mongoDb      *mgo.Database
	fileStore    Storage
}

var ()

// Config options
type Config struct {
	assetsPath     string
	templatesPath  string
	servingAddress string
	mongoAddress   string
	dataPath       string
	shutdownWait   time.Duration
}

var (
	// Config optoins TODO read from config file and/or params
	config = Config{
		assetsPath:     "./assets/",
		templatesPath:  "./templates/",
		dataPath:       "./data", // Mongo, files, indexes ...
		servingAddress: "127.0.0.1:5533",
		shutdownWait:   15 * time.Second,
		mongoAddress:   "127.0.0.1:27017",
	}

	entryServer = new(EntryServer)
)

func init() {
	// TODO use flag to initialize Config from command-line params and later to initialize from json/toml/other file-based config format.
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	entryServer.init(&config)
}

func (es *EntryServer) init(config *Config) (err error) {
	es.mongoSession, err = mgo.Dial(config.mongoAddress)
	if err != nil {
		return
	}
	es.mongoDb = es.mongoSession.DB(edraj)
	es.fileStore.RootPath = path.Join(config.dataPath, edraj)
	es.fileStore.TrashPath = path.Join(config.dataPath, trash, edraj)
	return
}

// Helper function
func entryObject(objectType string, e *Entry, createIfNil bool) (doc interface{}) {
	fieldName := strings.Title(objectType)
	field := reflect.ValueOf(e).Elem().FieldByName(fieldName)
	doc = field.Interface()
	if createIfNil && field.IsNil() /*&& field.Pointer() == 0*/ {
		field.Set(reflect.Indirect(reflect.New(field.Type().Elem())).Addr())
		doc = field.Interface()
	}
	return
}

// Create ...
func (es *EntryServer) Create(ctx context.Context, request *EntryRequest) (response *Receipt, err error) {
	response = &Receipt{}
	glog.Info(request)
	if request.Entry == nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = fmt.Errorf("Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}
	entryType := strings.ToLower(request.Entry.Type.String())

	doc := entryObject(entryType, request.Entry, false)

	err = es.mongoDb.C(entryType).Insert(doc)
	if err != nil {
		response.Status.Code = int32(codes.InvalidArgument)
		response.Status.Message = fmt.Sprintf("Failed to create entry (%v).", request.Entry)
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Update ...
func (es *EntryServer) Update(ctx context.Context, request *EntryRequest) (response *Receipt, err error) {
	response = &Receipt{}
	glog.Info(request)
	if request.Entry == nil /*|| request.Entry.Id == ""*/ {
		response.Status.Code = int32(codes.InvalidArgument)
		err = fmt.Errorf("Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.Entry.Type.String())

	doc := entryObject(entryType, request.Entry, false)

	err = es.mongoDb.C(entryType).Update(request.Entry.Id, doc)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to update entry (%v).", request.Entry)
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Query ...
func (es *EntryServer) Query(ctx context.Context, request *QueryRequest) (response *Response, err error) {
	response = &Response{}
	glog.Info(request)

	if request.Query == nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = fmt.Errorf("Query details are missing (%v).", request.Query)
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
	err = es.mongoDb.C(entryType).Find(query).All(objects.Interface())

	for i := 0; i < objects.Elem().Len(); i++ {
		entry := &Entry{}
		reflect.ValueOf(entry).Elem().FieldByName(fieldName).Set(reflect.ValueOf(objects.Elem().Index(i).Addr().Interface()))
		response.Entries = append(response.Entries, entry)
	}

	if err != nil {
		response.Status.Code = int32(codes.Internal)
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Get ...
func (es *EntryServer) Get(ctx context.Context, request *IdRequest) (response *Response, err error) {
	glog.Info(request)
	response = &Response{}

	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = fmt.Errorf("EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.EntryType.String())
	entry := Entry{}
	doc := entryObject(entryType, &entry, true)

	err = es.mongoDb.C(entryType).FindId(request.EntryId /*bson.ObjectIdHex(request.EntryID)*/).One(doc)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to get entry (%v).", request.EntryId)
		return
	}

	response.Status.Code = int32(codes.OK)
	response.Entries = []*Entry{&entry}
	response.Returned = 1
	response.Total = 1

	return
}

// Delete ...
func (es *EntryServer) Delete(ctx context.Context, request *IdRequest) (response *Receipt, err error) {
	glog.Info(request)
	response = &Receipt{}
	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = fmt.Errorf("EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.EntryType.String())
	err = es.mongoDb.C(entryType).Remove(&struct{ _id string }{_id: request.EntryId})
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to delete entry (%v).", request.EntryId)
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Run ...
func Run() error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterEntryServiceServer(server, entryServer)
	glog.Info("Starting the gRPC server")
	server.Serve(listen)
	return nil
}

func main() {
	defer glog.Flush()

	if err := Run(); err != nil {
		glog.Fatal(err)
	}
}
