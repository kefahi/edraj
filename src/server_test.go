package main

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

func checkSuccess(t *testing.T, response interface{}, err error) {
	if err != nil {
		t.Errorf("HelloTest(%v) got unexpected error", err)
	}
}

func checkFailure(t *testing.T, response interface{}, err error) {
	if err == nil {
		switch r := response.(type) {
		case *Response:
			t.Errorf("Response: %v", r)
		case *Receipt:
			t.Errorf("Receipt: %v", r)
		}
	}

}

func TestEntryGRPC(t *testing.T) {

	fmt.Println("Testing gRPC server")

	one := Content{Id: "one", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}
	two := Content{Id: "two", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}

	var receipt *Receipt
	var response *Response
	var err error

	receipt, err = entryGrpc.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id})
	//checkFailure(t, receipt, err)
	receipt, err = entryGrpc.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id})
	//checkFailure(t, receipt, err)

	receipt, err = entryGrpc.Create(context.Background(), &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &one}})
	checkSuccess(t, receipt, err)
	receipt, err = entryGrpc.Create(context.Background(), &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &two}})
	checkSuccess(t, receipt, err)

	response, err = entryGrpc.Query(context.Background(), &QueryRequest{Query: &Query{EntryType: EntryType_CONTENT}})
	checkSuccess(t, response, err)
	response, err = entryGrpc.Get(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id})
	checkSuccess(t, response, err)

	receipt, err = entryGrpc.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id})
	checkSuccess(t, receipt, err)
	receipt, err = entryGrpc.Delete(context.Background(), &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id})
	checkSuccess(t, receipt, err)

}
