package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	mgo "gopkg.in/mgo.v2"
)

/*
//import
//	"github.com/jban332/kin-openapi/openapi3"
//	"github.com/jban332/kin-openapi/openapi3filter"


func getOperation(req *http.Request) (*openapi3.Operation, error) {
	// Load Swagger file
	router := openapi3filter.NewRouter().WithSwaggerFromFile("swagger.json")

	// Find route
	route, _, err := router.FindRoute("GET", req.URL)
	if err != nil {
		return nil, err
	}

	// Get OpenAPI 3 operation
	return route.Operation, nil
}

var router = openapi3filter.NewRouter().WithSwaggerFromFile("swagger.json")

func validateRequest(req *http.Request) {
	openapi3filter.ValidateRequest(nil, &openapi3filter.RequestValidationInput{
		Request: req,
		//Route:   router,
	})

	// Get response

	openapi3filter.ValidateResponse(nil, &openapi3filter.ResponseValidationInput{
	// ...
	})
}
*/

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
	// Config optoins, currently all hardcoded. TODO read from config file and/or params
	config = Config{
		assetsPath:     "./assets/",
		templatesPath:  "./templates/",
		dataPath:       "./data", // Mongo, files, indexes ...
		servingAddress: "127.0.0.1:5533",
		shutdownWait:   15 * time.Second,
		mongoAddress:   "127.0.0.1:27017",
	}

	mongoSession *mgo.Session
	entryService EntryService
)

func init() {
	// TODO use flag to initialize Config from command-line params and later to initialize from json/toml/other file-based config format.
	var err error
	mongoSession, err = mgo.Dial(config.mongoAddress)
	if err != nil {
		log.Fatal(err)
	}

	entryService.init(config)
}

// Log wrapper for all http calls
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("index.html")).Execute(w, "Hello World!")
}

func main() {

	http.HandleFunc("/api/entry/", EntryAPI)
	http.HandleFunc("/hello", HelloAPI)
	http.Handle("/assets", http.FileServer(http.Dir(config.assetsPath)))
	http.HandleFunc("/", index)
	r := Log(http.DefaultServeMux)

	server := &http.Server{
		Addr: config.servingAddress,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Started server at http://%s\n", config.servingAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

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
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
