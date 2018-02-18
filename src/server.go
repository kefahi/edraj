package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

// Config options
type Config struct {
	assetsPath     string
	certsPath      string
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
		certsPath:      "../certs/",
		templatesPath:  "./templates/",
		dataPath:       "./data", // Mongo, files, indexes ...
		servingAddress: "127.0.0.1:5533",
		shutdownWait:   15 * time.Second,
		mongoAddress:   "127.0.0.1:27017",
	}

	entryGrpc = new(EntryGRPC)
	entryMan  = &EntryMan{}

	httpServer *http.Server
	grpcServer *grpc.Server

	grpc2http = map[codes.Code]int{
		codes.OK: http.StatusOK,
	}
)

func init() {
	// TODO use flag to initialize Config from command-line params and later to initialize from json/toml/other file-based config format.
	flag.Parse()
	// for glog flag.Lookup("logtostderr").Value.Set("true")
	grpc.EnableTracing = true
	grpclog.SetLogger(log.New(os.Stdout, "edraj: ", log.LstdFlags))
	entryMan.init(&config)
}

// TODO enable Tracing and test it

func main() {
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
	grpclog.Info("shutting down")
	os.Exit(0)
}
