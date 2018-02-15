package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"os/signal"

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
	trash = "trash"
	edraj = "edraj"
)

// EntryServer implements EntryServiceServer
type EntryServer struct {
	mongoSession *mgo.Session
	mongoDb      *mgo.Database
	fileStore    Storage
}

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
	httpServer  *http.Server
	grpcServer  *grpc.Server
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
	if createIfNil && field.IsNil() {
		field.Set(reflect.Indirect(reflect.New(field.Type().Elem())).Addr())
		doc = field.Interface()
	}
	return
}

// Create ...
func (es *EntryServer) Create(ctx context.Context, request *EntryRequest) (response *Receipt, err error) {
	response = &Receipt{Status: &Status{}}
	glog.Info(request)
	if request.Entry == nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = grpc.Errorf(codes.InvalidArgument, "Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}
	entryType := strings.ToLower(request.Entry.Type.String())

	doc := entryObject(entryType, request.Entry, false)

	err = es.mongoDb.C(entryType).Insert(doc)
	if err != nil {
		response.Status.Code = int32(codes.Internal)
		response.Status.Message = fmt.Sprintf("Failed to create entry (%v).", request.Entry)
		err = grpc.Errorf(codes.Internal, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Update ...
func (es *EntryServer) Update(ctx context.Context, request *EntryRequest) (response *Receipt, err error) {
	response = &Receipt{Status: &Status{}}
	glog.Info(request)
	if request.Entry == nil || request.Entry.Id == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = grpc.Errorf(codes.Internal, "Entry details are missing (%v).", request.Entry)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.Entry.Type.String())

	doc := entryObject(entryType, request.Entry, false)

	err = es.mongoDb.C(entryType).Update(request.Entry.Id, doc)
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to update entry (%v).", request.Entry)
		err = grpc.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Query ...
func (es *EntryServer) Query(ctx context.Context, request *QueryRequest) (response *Response, err error) {
	response = &Response{Status: &Status{}}
	glog.Info(request)

	if request.Query == nil {
		response.Status.Code = int32(codes.InvalidArgument)
		err = grpc.Errorf(codes.InvalidArgument, "Query details are missing (%v).", request.Query)
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
		err = grpc.Errorf(codes.Internal, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Get ...
func (es *EntryServer) Get(ctx context.Context, request *IdRequest) (response *Response, err error) {
	glog.Info(request)
	response = &Response{Status: &Status{}}

	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = grpc.Errorf(codes.InvalidArgument, "EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
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
		err = grpc.Errorf(codes.NotFound, err.Error())
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
	response = &Receipt{Status: &Status{}}
	if request.EntryId == "" {
		response.Status.Code = int32(codes.InvalidArgument)
		err = grpc.Errorf(codes.InvalidArgument, "EntryId (%v) or EntryType (%v) are missing.", request.EntryId, request.EntryType)
		response.Status.Message = err.Error()
		return
	}

	entryType := strings.ToLower(request.EntryType.String())
	err = es.mongoDb.C(entryType).Remove(&struct{ _id string }{_id: request.EntryId})
	if err != nil {
		response.Status.Code = int32(codes.NotFound)
		response.Status.Message = fmt.Sprintf("Failed to delete entry (%v).", request.EntryId)
		err = grpc.Errorf(codes.NotFound, err.Error())
		return
	}

	response.Status.Code = int32(codes.OK)
	return
}

// Log wrapper for all http calls
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Run ...
func runGRPC() {
	listen, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		glog.Fatal(err)
	}
	grpcServer = grpc.NewServer()
	RegisterEntryServiceServer(grpcServer, entryServer)
	glog.Info("Starting the gRPC server")
	grpcServer.Serve(listen)
}

func index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("index.html")).Execute(w, "Hello World!")
}
func runHTTP() {
	http.Handle("/assets", http.FileServer(http.Dir(config.assetsPath)))
	http.HandleFunc("/", index)
	r := Log(http.DefaultServeMux)

	httpServer = &http.Server{
		Addr: config.servingAddress,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	glog.Infof("Started server at http://%s\n", config.servingAddress)
	if err := httpServer.ListenAndServe(); err != nil {
		glog.Fatal(err)
	}
}

func mainServer() {
	go runGRPC()
	go runHTTP()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), config.shutdownWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	httpServer.Shutdown(ctx)
	grpcServer.Stop()
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	glog.Info("shutting down")
	os.Exit(0)
}

func main() {
	defer glog.Flush()

	command := mainServer
	if len(os.Args) > 1 && os.Args[1] == "client" {
		command = mainClient
	}
	command()
}
