# Golang go-libp2p-pubsub

Study of new go-libp2p for introduction protobuf with protocol.

Stack of development

* go 1.13+
* go-libp2p-pubsub 0.2.0
* protobuf 1.3.0

### Install dependence

Before start you need install some deps. Follow the sequence.

```shell=
$ go get -u -v github.com/golang/protobuf/proto
$ go get -u -v github.com/golang/protobuf/protoc-gen-go
```

If you have a **Linux** distribution you need install `protoc`

```shell=
$ sudo apt install protobuf-compiler
```

### Generate the protocol

Just you need run the command for compile you `.proto` on yours folder

```shell=
$ protoc --go_out=. *.proto
```
