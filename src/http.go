package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

// Log wrapper for all http calls
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		grpclog.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("index.html")).Execute(w, "Hello World!")
}

// HelloAPI ...
func HelloAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func respond(w http.ResponseWriter, response *Response) {
	data, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CODE", fmt.Sprintf("(%d): %s", response.Status.Code, response.Status.Message))
	w.WriteHeader(grpc2http[codes.Code(response.Status.Code)])
	w.Write(data)
}

// EntryAPI serves the api
// TODO validate requester id (by checking the signature)
// TODO validate access-control
func EntryAPI(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var err error

	if strings.HasPrefix(r.RequestURI, "/api/entry/query") {
		var request QueryRequest
		if err = json.Unmarshal(body, &request); err == nil {
			response, _ := entryGrpc.Query(context.Background(), &request)
			respond(w, response)
			return
		}
	} else if strings.HasPrefix(r.RequestURI, "/api/entry/get") {
		var request IdRequest
		if err := json.Unmarshal(body, &request); err == nil {
			response, _ := entryGrpc.Get(context.Background(), &request)
			respond(w, response)
			return
		}
	}

	respond(w, &Response{
		Status: &Status{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		}})
}

func runHTTP() {

	http.HandleFunc("/api/entry/", EntryAPI)
	http.HandleFunc("/api/hello", HelloAPI)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(config.assetsPath))))
	//http.HandleFunc("/", index)
	r := Log(http.DefaultServeMux)

	httpServer = &http.Server{
		Addr: config.servingAddress,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	//grpclog.Infof("Started server at http://%s\n", config.servingAddress)
	grpclog.Infof("Starting server. http://%s", config.servingAddress)
	if err := httpServer.ListenAndServe(); err != nil {
		grpclog.Fatal(err)
	}
}
