```
protoc.exe --go_out=.\grpc --go-grpc_out=.\grpc .\grpc\example.proto
```

```
 cd cd .\server\
 go run .\server.go
```
```
 cd cd .\client\
  go run .\client.go
```