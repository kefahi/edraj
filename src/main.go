package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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
	// Config optoins, currently all hardcoded. TODO read from config file and/or params
	config = Config{
		assetsPath:     "./assets/",
		templatesPath:  "./templates/",
		dataPath:       "./data", // Mongo, files, indexes ...
		servingAddress: "127.0.0.1:5533",
		shutdownWait:   15 * time.Second,
		mongoAddress:   "127.0.0.1:27017",
	}

	entryService EntryService
)

func init() {
	// TODO use flag to initialize Config from command-line params and later to initialize from json/toml/other file-based config format.
	entryService.init(config)

}

// Log wrapper for all http calls
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {

	/*
		r := mux.NewRouter()
		r.HandleFunc("/api/entry", QueryAPI).Methods("GET")
		r.HandleFunc("/api/entry/{id}", GetEntryAPI).Methods("GET")
		r.HandleFunc("/api/entry", CreateEntryAPI).Methods("PUT")
		r.HandleFunc("/api/entry/{id}", UpdateEntryAPI).Methods("POST")
		r.HandleFunc("/api/entry/{id}", DeleteEntryAPI).Methods("DELETE")
		r.HandleFunc("/api/domains", NotImplementedAPI).Methods("GET")
		r.HandleFunc("/api/hello", HelloAPI).Methods("GET")
		r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(config.assetsPath))))
	*/

	http.HandleFunc("/api/entry/", EntryAPI)
	http.HandleFunc("/hello", HelloAPI)
	http.Handle("/assets", http.FileServer(http.Dir(config.assetsPath)))
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
