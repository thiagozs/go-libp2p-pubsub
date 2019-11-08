package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	router "github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
	core "github.com/thiagozs/go-libp2p-pubsub/core"
	ds "github.com/thiagozs/go-libp2p-pubsub/mdns"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// chainoptions for transport
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	// chainoptions for muxers
	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	// security instance
	security := libp2p.Security(secio.ID, secio.New)

	// listenAddress
	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)

	// create a hash table addrs
	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (router.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		return dht, err
	}
	// routing with opt with dht
	routing := libp2p.Routing(newDHT)

	// start libp2p
	host, err := libp2p.New(ctx, transports, listenAddrs, muxers, security, routing)
	if err != nil {
		panic(err)
	}

	// start pubsub engine
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	// subscribe for a topic
	sub, err := ps.Subscribe(core.Topic)
	if err != nil {
		panic(err)
	}

	// Start counters
	cc := core.NewCounters()

	// handler for messages
	go core.PubsubHandler(cc, ctx, sub)

	// list address for listering
	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	// change for real bootstrap
	targetAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/4001/ipfs/QmXBbs4x3E7f8TNgZPCcG38vgRFiZdeAC6qPPVHKGtAR7x")
	if err != nil {
		panic(err)
	}

	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		panic(err)
	}

	// connect for the host
	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to", targetInfo.ID)

	// start domain name server
	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		panic(err)
	}
	mdns.RegisterNotifee(&ds.Notifee{H: host, CTX: ctx})

	// tells the DHT to get into a bootstrapped state satisfying the IpfsRouter interface
	err = dht.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	// handler chat input
	donec := make(chan struct{}, 1)
	go core.ChatInputLoop(cc, ctx, host, ps, donec)

	// monitoring program stats
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop: // exit program
		dht.Close()
		host.Network().Close()
		host.Close()
		os.Exit(0)
	case <-donec: // finish command close with host
		host.Close()
	}
}
