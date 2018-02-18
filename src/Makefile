
.PHONY: all
all: proto edraj ../certs

../certs: 
	../bin/certstrap --depot-path=../certs init --passphrase "" -o edraj -cn edrajRootCA
	../bin/certstrap --depot-path=../certs request-cert --passphrase "" --domain edraj.io
	../bin/certstrap --depot-path=../certs sign --CA edrajRootCA edraj.io
	../bin/certstrap --depot-path=../certs request-cert --passphrase "" --domain localhost
	../bin/certstrap --depot-path=../certs sign --CA edrajRootCA localhost
	../bin/certstrap --depot-path=../certs request-cert --passphrase "" -cn kefah
	../bin/certstrap --depot-path=../certs sign --CA edrajRootCA kefah

.PHONY: proto
proto: edraj.pb.go ./python_client/edraj_pb2_grpc.py

edraj.pb.go: ../api-spec/edraj.proto 
	protoc -I ../api-spec/ ../api-spec/edraj.proto --go_out=plugins=grpc:. 
	sed -i 's/json:"\([^,"]*\),omitempty"/json:"\1,omitempty" bson:"\1,omitempty"/g' ./edraj.pb.go
	sed -i 's/bson:"id,omitempty"/bson:"_id,omitempty"/g' ./edraj.pb.go 
	sed -i '1{s/$$/  DO NOT EDIT./}' ./edraj.pb.go

./python_client/edraj_pb2_grpc.py: ../api-spec/edraj.proto
	python3 -m grpc_tools.protoc --proto_path=../api-spec --python_out=./python_client/ --grpc_python_out=./python_client/ ../api-spec/edraj.proto


edraj: server.go edraj.pb.go storage.go entry.go
	go build -o edraj server.go edraj.pb.go storage.go entry.go grpc.go http.go

.PHONY: test
test: 
	go test

.PHONY: clean
clean: 
	rm -f edraj edraj.pb.go ./python_client/edraj_pb2_grpc.py ./python_client/edraj_pb2.py
