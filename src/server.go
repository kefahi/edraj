package main

import (
	"flag"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

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

	entryGrpc = new(EntryGRPC)

	httpServer *http.Server
	grpcServer *grpc.Server
)

func init() {
	// TODO use flag to initialize Config from command-line params and later to initialize from json/toml/other file-based config format.
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	entryGrpc.entryMan = &EntryMan{}
	entryGrpc.entryMan.init(&config)
}

// EntryGRPC implements EntryServiceServer and delegates the calls to EntryMan
type EntryGRPC struct {
	entryMan *EntryMan
}

// Create ...
func (es *EntryGRPC) Create(ctx context.Context, request *EntryRequest) (*Receipt, error) {
	return es.entryMan.create(request)
}

// Update ...
func (es *EntryGRPC) Update(ctx context.Context, request *EntryRequest) (*Receipt, error) {
	return es.entryMan.update(request)

}

// Query ...
func (es *EntryGRPC) Query(ctx context.Context, request *QueryRequest) (*Response, error) {
	return es.entryMan.query(request)
}

// Get ...
func (es *EntryGRPC) Get(ctx context.Context, request *IdRequest) (*Response, error) {
	return es.entryMan.get(request)
}

// Delete ...
func (es *EntryGRPC) Delete(ctx context.Context, request *IdRequest) (*Receipt, error) {
	return es.entryMan.delete(request)
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
	RegisterEntryServiceServer(grpcServer, entryGrpc)
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

func main() {
	defer glog.Flush()
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

/*
func main() {
	command := mainServer
	if len(os.Args) > 1 && os.Args[1] == "client" {
		command = mainClient
	}
	command()
}*/