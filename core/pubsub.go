package core

import (
	"context"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	buffer "github.com/thiagozs/go-libp2p-pubsub/proto/v2"
)

var handles = map[string]string{}

// Topic channel
const Topic = "/libp2p-pubsub/chat/thiagozs"

func pubsubMessageHandler(id peer.ID, msg *buffer.SendMessage) {
	handle, ok := handles[id.String()]
	if !ok {
		handle = id.ShortString()
	}
	fmt.Printf("%s: %s\n", handle, msg.Data)
}

func pubsubUpdateHandler(id peer.ID, msg *buffer.UpdatePeer) {
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

		req := &buffer.Request{}
		err = proto.Unmarshal(msg.Data, req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		switch req.Type.String() {
		case buffer.Request_SEND_MESSAGE.String():
			pubsubMessageHandler(msg.GetFrom(), req.SendMessage)
		case buffer.Request_UPDATE_PEER.String():
			pubsubUpdateHandler(msg.GetFrom(), req.UpdatePeer)
		}
	}
}
