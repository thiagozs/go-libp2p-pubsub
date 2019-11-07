package core

import (
	"context"
	"fmt"
	"os"

	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	proto "github.com/thiagozs/go-libp2p-pubsub/proto/v1"
)

var handles = map[string]string{}

// Topic channel
const Topic = "/libp2p-pubsub/chat/thiagozs"

func pubsubMessageHandler(id peer.ID, msg *proto.SendMessage) {
	handle, ok := handles[id.String()]
	if !ok {
		handle = id.ShortString()
	}
	fmt.Printf("%s: %s\n", handle, msg.Data)
}

func pubsubUpdateHandler(id peer.ID, msg *proto.UpdatePeer) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = string(msg.UserHandle)
	fmt.Printf("%s -> %s\n", oldHandle, msg.UserHandle)
}

// PubsubHandler start listenner and send message
func PubsubHandler(ctx context.Context, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		req := &proto.Request{}
		err = req.Unmarshal(msg.Data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		switch *req.Type {
		case proto.Request_SEND_MESSAGE:
			pubsubMessageHandler(msg.GetFrom(), req.SendMessage)
		case proto.Request_UPDATE_PEER:
			pubsubUpdateHandler(msg.GetFrom(), req.UpdatePeer)
		}
	}
}
