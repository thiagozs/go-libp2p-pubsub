package core

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	buffer "github.com/thiagozs/go-libp2p-pubsub/proto/v2"
)

func sendMessage(ps *pubsub.PubSub, msg string) {
	msgId := make([]byte, 10)
	_, err := rand.Read(msgId)
	defer func() {
		if err != nil {
			fmt.Printf("Erro on sendMessage : %s\n", err)
		}
	}()
	if err != nil {
		return
	}
	now := time.Now().Unix()
	req := &buffer.Request{
		Type: buffer.Request_SEND_MESSAGE,
		SendMessage: &buffer.SendMessage{
			Id:      msgId,
			Data:    []byte(msg),
			Created: now,
		},
	}
	msgBytes, err := proto.Marshal(req)
	if err != nil {
		fmt.Printf("Error on marshal message : %s\n", err)
		return
	}

	err = ps.Publish(Topic, msgBytes)
	if err != nil {
		fmt.Printf("Error on publish message : %s\n", err)
		return
	}
}

func updatePeer(ps *pubsub.PubSub, id peer.ID, handle string) {
	oldHandle, ok := handles[id.String()]
	if !ok {
		oldHandle = id.ShortString()
	}
	handles[id.String()] = handle

	req := &buffer.Request{
		Type: buffer.Request_UPDATE_PEER,
		UpdatePeer: &buffer.UpdatePeer{
			UserHandle: []byte(handle),
		},
	}
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		fmt.Printf("Error on marshal message : %s\n", err)
		return
	}
	err = ps.Publish(Topic, reqBytes)
	if err != nil {
		fmt.Printf("Error on publish message : %s\n", err)
		return
	}

	fmt.Printf("%s -> %s\n", oldHandle, handle)
}
