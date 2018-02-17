package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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

	httpServer *http.Server
	grpcServer *grpc.Server
)

func init() {
	// TODO use flag to initialize Config from command-line params and later to initialize from json/toml/other file-based config format.
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	grpc.EnableTracing = true
	entryGrpc.entryMan = &EntryMan{}
	entryGrpc.entryMan.init(&config)
}

// EntryGRPC implements EntryServiceServer and delegates the calls to EntryMan
type EntryGRPC struct {
	entryMan *EntryMan
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()
	//newStream := grpc_middleware.WrapServerStream(stream)
	//newStream.WrappedContext = context.WithValue(ctx, "user_id", "john@example.com")
	err = handler(srv, stream)
	glog.Infof("invoke stream method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return
}

//grpc.StreamServerInterceptor()
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		// Read edraj-signature-bin if invalid:
		// return nil, grpc.Errorf(codes.Unauthenticated, "invalid signature")
		for k, v := range headers {
			if strings.HasPrefix(k, "edraj-") {
				glog.Infof("Ctx %v: %v", k, v)
			}
		}
	} else {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	grpc.SendHeader(ctx, metadata.New(map[string]string{"edraj-header": "my-value"}))
	grpc.SetTrailer(ctx, metadata.New(map[string]string{"edraj-trailer": "my-value"}))
	//ctx = context.WithValue(ctx, "user_id", "john@example.com")
	start := time.Now()
	resp, err := handler(ctx, req)
	glog.Infof("invoke server method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

// TODO serverStreamInterceptor

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

func logSleep(ctx context.Context, d time.Duration) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("sleeping for %s", d)
	}
}

// Notifications ...
func (es *EntryGRPC) Notifications(request *QueryRequest, stream EntryService_NotificationsServer) (err error) {
	// TODO establish per-call (user/call notification channel)
	// TODO handle cancelation
	ctx := stream.Context()
	for i := 0; i < 10; i++ {
		d := time.Duration(rand.Intn(10)) * time.Second
		logSleep(ctx, d)
		select {
		case <-time.After(d):
			err := stream.Send(&Notification{
				What:      fmt.Sprintf("result %d for [%s] from backend %d", i, request, d),
				Timestamp: uint64(i),
			})
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	/*
		for j := 0; j < 10; j++ {

			if err := stream.Send(&Notification{}); err != nil {
				return err
			}
		}

	*/
	return nil

}

// TODO enable Tracing and test it

// Log wrapper for all http calls
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Run ...
func runGRPC() {
	listen, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		glog.Fatal(err)
	}
	certificate, err := tls.LoadX509KeyPair(path.Join(config.certsPath, "localhost.crt"), path.Join(config.certsPath, "localhost.key"))

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(path.Join(config.certsPath, "edrajRootCA.crt"))
	if err != nil {
		glog.Fatalf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		glog.Fatal("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	// TODO : Additionally consider grpc.StatsHandler(th)
	grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.MaxConcurrentStreams(64),
		//grpc.InTapHandle(NewTap.Handler),
	)

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
