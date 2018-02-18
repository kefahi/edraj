package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

// EntryGRPC implements EntryServiceServer and delegates the calls to EntryMan
type EntryGRPC struct {
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()
	//newStream := grpc_middleware.WrapServerStream(stream)
	//newStream.WrappedContext = context.WithValue(ctx, "user_id", "john@example.com")
	err = handler(srv, stream)
	grpclog.Infof("invoke stream method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return
}

//grpc.StreamServerInterceptor()
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		// Read edraj-signature-bin if invalid:
		// return nil, grpc.Errorf(codes.Unauthenticated, "invalid signature")
		for k, v := range headers {
			if strings.HasPrefix(k, "edraj-") {
				grpclog.Infof("Ctx %v: %v", k, v)
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
	grpclog.Infof("invoke unary method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

// TODO serverStreamInterceptor

// Create ...
func (es *EntryGRPC) Create(ctx context.Context, request *EntryRequest) (*Receipt, error) {
	return entryMan.create(request)
}

// Update ...
func (es *EntryGRPC) Update(ctx context.Context, request *EntryRequest) (*Receipt, error) {
	return entryMan.update(request)
}

// Query ...
func (es *EntryGRPC) Query(ctx context.Context, request *QueryRequest) (*Response, error) {
	return entryMan.query(request)
}

// Get ...
func (es *EntryGRPC) Get(ctx context.Context, request *IdRequest) (*Response, error) {
	return entryMan.get(request)
}

// Delete ...
func (es *EntryGRPC) Delete(ctx context.Context, request *IdRequest) (*Receipt, error) {
	return entryMan.delete(request)
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
		d := time.Duration(rand.Intn(3)) * time.Second
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

// Run ...
func runGRPC() {
	listen, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		grpclog.Fatal(err)
	}
	certificate, err := tls.LoadX509KeyPair(path.Join(config.certsPath, "localhost.crt"), path.Join(config.certsPath, "localhost.key"))

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(path.Join(config.certsPath, "edrajRootCA.crt"))
	if err != nil {
		grpclog.Fatalf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		grpclog.Fatal("failed to append client certs")
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
	grpclog.Info("Starting the gRPC server")
	grpcServer.Serve(listen)
}
