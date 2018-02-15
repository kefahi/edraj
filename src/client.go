package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// EntryClient ...
type EntryClient struct {
	address string
	conn    *grpc.ClientConn
	service EntryServiceClient
}

func (ec *EntryClient) init() {
	var err error
	ec.conn, err = grpc.Dial(ec.address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ec.service = NewEntryServiceClient(ec.conn)
}

func (ec *EntryClient) close() {
	ec.conn.Close()
}

func check(response interface{}, err error) {
	if err != nil {
		log.Printf("call Failed: %v", err)
	} else {
		switch r := response.(type) {
		case *Response:
			log.Printf("Response: %v", r)
		case *Receipt:
			log.Printf("Receipt: %v", r)
		}
	}
}

func mainClient() {

	// Set up a connection to the server.
	client := EntryClient{address: "localhost:50051"}
	client.init()
	defer client.close()

	one := Content{Id: "one", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}
	two := Content{Id: "two", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}

	check(client.service.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}))
	check(client.service.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id}))

	check(client.service.Create(context.Background(), &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &one}}))
	check(client.service.Create(context.Background(), &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &two}}))

	check(client.service.Query(context.Background(), &QueryRequest{Query: &Query{EntryType: EntryType_CONTENT}}))
	check(client.service.Get(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}))

	check(client.service.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}))
	check(client.service.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id}))
}
