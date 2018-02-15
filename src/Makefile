all: gserver gclient

edraj.pb.go: ../api-spec/edraj.proto 
	protoc -I ../api-spec/ ../api-spec/edraj.proto --go_out=plugins=grpc:. && sed -i 's/ json:"id,omitempty"/ json:"id,omitempty" bson:"_id"/g' ./edraj.pb.go

gserver: server.go edraj.pb.go storage.go
	go build -o gserver server.go edraj.pb.go storage.go

gclient: client.go edraj.pb.go
	go build -o gclient client.go edraj.pb.go

server: gserver
	./gserver

client: gclient
	./gclient

clean: 
	rm -f gserver gclient edraj.pb.go
