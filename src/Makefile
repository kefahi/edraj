
.PHONY: all
all: edraj ../certs

../certs: 
	../bin/certstrap --depot-path=../certs init --passphrase "" -o edraj -cn edrajRootCA
	../bin/certstrap --depot-path=../certs request-cert --passphrase "" --domain edraj.io
	../bin/certstrap --depot-path=../certs sign --CA edrajRootCA edraj.io
	../bin/certstrap --depot-path=../certs request-cert --passphrase "" --domain localhost
	../bin/certstrap --depot-path=../certs sign --CA edrajRootCA localhost
	../bin/certstrap --depot-path=../certs request-cert --passphrase "" -cn kefah
	../bin/certstrap --depot-path=../certs sign --CA edrajRootCA kefah


edraj.pb.go: ../api-spec/edraj.proto 
	protoc -I ../api-spec/ ../api-spec/edraj.proto --go_out=plugins=grpc:. && sed -i 's/ json:"id,omitempty"/ json:"id,omitempty" bson:"_id"/g' ./edraj.pb.go && sed -i '1{s/$$/  DO NOT EDIT./}' ./edraj.pb.go

edraj: server.go edraj.pb.go storage.go entry.go
	go build -o edraj server.go edraj.pb.go storage.go entry.go

.PHONY: test
test: 
	go test

.PHONY: clean
clean: 
	rm -f edraj edraj.pb.go
