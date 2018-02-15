## Installation
`go get github.com/BurntSushi/toml gopkg.in/mgo.v2 github.com/gorilla/mux`

```
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

## Go Meta Linter
```
go get -u github.com/alecthomas/gometalinter
gometalinter --install
gometalinter ./ --exclude='main redeclared' --exclude='edraj.pb.go' --exclude='other declaration of main'
```


`dlv debug --headless --listen=:2345 --log -- -myArg=123`

## TODO
+ Fix Delete method
+ Add more generic tests

