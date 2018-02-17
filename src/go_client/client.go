package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// EntryClient ...
type EntryClient struct {
	host      string
	port      int
	certsPath string
	conn      *grpc.ClientConn
	service   EntryServiceClient
}

func (ec *EntryClient) init() {
	var err error
	grpc.EnableTracing = true
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
		ServerName:   ec.host,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	// TODO : additionally consider stats: grpc.WithStatsHandler(th)
	ec.conn, err = grpc.Dial(
		fmt.Sprintf("%s:%d", ec.host, ec.port),
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithUnaryInterceptor(clientUnaryInterceptor),
		grpc.WithStreamInterceptor(clientStreamInterceptor))

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

func clientUnaryInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	headers := metadata.New(
		map[string]string{
			"edraj-signature-bin": "mysig",
			"edraj-pubkey-bin":    "mykey",
			"edraj-id":            "myid",
			"edraj-timestamp":     "mytime"})

	//log.Println(headers)

	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
	// this is the critical step that includes your headers
	ctx = metadata.NewOutgoingContext(ctx, headers)
	var header, trailer metadata.MD

	opts = append(opts, grpc.Header(&header))
	opts = append(opts, grpc.Trailer(&trailer))

	err := invoker(ctx, method, req, reply, cc, opts...) // <==
	printReturnedMeta(header, trailer)
	log.Printf("invoke remote method=%s duration=%s error=%v", method, time.Since(start), err)
	return err
}

func clientStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (client grpc.ClientStream, err error) {
	start := time.Now()
	client, err = streamer(ctx, desc, cc, method, opts...)
	log.Printf("invoke remote stream method=%s duration=%s error=%v", method, time.Since(start), err)
	return
}

// TODO clientStreamInterceptor

func main() {

	// Set up a connection to the server.
	client := EntryClient{host: "localhost", port: 50051, certsPath: "../../certs/"}
	client.init()
	defer client.close()

	ctx := context.Background()
	one := Content{Id: "one", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}
	two := Content{Id: "two", Path: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}

	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}))

	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id}))

	check(client.service.Create(ctx, &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &one}}))
	check(client.service.Create(ctx, &EntryRequest{Entry: &Entry{Type: EntryType_CONTENT, Content: &two}}))

	check(client.service.Query(ctx, &QueryRequest{Query: &Query{EntryType: EntryType_CONTENT}}))
	check(client.service.Get(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}))

	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: one.Id}))
	check(client.service.Delete(ctx, &IdRequest{EntryType: EntryType_CONTENT, EntryId: two.Id}))

	stream, err := client.service.Notifications(ctx, &QueryRequest{})
	if err != nil {
		log.Println("Error on streaming", err)
		return
	}
	for {
		notification, err := stream.Recv()
		if err == io.EOF {
			log.Println("Notifications stream ends here")
			break
		}
		if err != nil {
			log.Fatalf("%v.Notifications(_) = _, %v", client, err)
		}
		log.Println(notification)
	}
}
