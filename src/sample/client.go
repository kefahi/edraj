package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// EntryClient ...
type EntryClient struct {
	address   string
	certsPath string
	conn      *grpc.ClientConn
	service   EntryServiceClient
}

func (ec *EntryClient) init() {
	var err error

	certificate, err := tls.LoadX509KeyPair(
		path.Join(ec.certsPath, "kefah.crt"),
		path.Join(ec.certsPath, "kefah.key"),
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(path.Join(ec.certsPath, "edrajRootCA.crt"))
	if err != nil {
		log.Fatalf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	ec.conn, err = grpc.Dial(ec.address, grpc.WithTransportCredentials(transportCreds))

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

func printReturnedMeta(meta ...metadata.MD) {
	for i, one := range meta {
		for key, value := range one {
			fmt.Printf("[%d] %s => %s\n", i, key, value)
		}
	}
}

func main() {

	// Set up a connection to the server.
	client := EntryClient{address: "localhost:50051", certsPath: "../../out/"}
	client.init()
	defer client.close()

	headers := metadata.New(map[string]string{"edraj-signature-bin": "mysig", "edraj-pubkey-bin": "mykey", "edraj-id": "myid", "edraj-timestamp": "mytime"})

	log.Println(headers)

	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
	// this is the critical step that includes your headers
	ctx := metadata.NewOutgoingContext(context.Background(), headers)
	var header, trailer metadata.MD

	one := Content{Id: "one", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}
	two := Content{Id: "two", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}

	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)
	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)

	check(client.service.Create(ctx, &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &one}}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)
	check(client.service.Create(ctx, &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &two}}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)

	check(client.service.Query(ctx, &QueryRequest{Query: &Query{EntryType: EntryType_CONTENT}}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)
	check(client.service.Get(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)

	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)
	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id}, grpc.Header(&header), grpc.Trailer(&trailer)))
	printReturnedMeta(header, trailer)
}
