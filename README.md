# Golang go-libp2p-pubsub

Study of new go-libp2p for introduction protobuf with protocol. Proof of concept.

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

### Chat commands

For test message stress, you can run a command for each peer. Just type

* `/msgtimer 100 100 metrics`, send a message with 100 miliseconds of delay, write a 100 messages with the slug `metrics`
* `/name your-nickname`, change our nickname
* `/stats`, show stats about the messages stress
* `/reset`, reset stats about the messages stress

### Compile

* `make build` or only `make`

### Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2019**, thiagozs.
