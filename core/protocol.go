package core

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"time"

	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	buffer "github.com/thiagozs/go-libp2p-pubsub/proto/v2"
)

func sendMessage(ps *pubsub.PubSub, msg string) {
	msgId := make([]byte, 10)
	_, err := rand.Read(msgId)
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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

// ChatInputLoop handler for message between peers
func ChatInputLoop(ctx context.Context, h host.Host, ps *pubsub.PubSub, donec chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()

		// change name
		if strings.HasPrefix(msg, "/name ") {
			newHandle := strings.TrimPrefix(msg, "/name ")
			newHandle = strings.TrimSpace(newHandle)
			updatePeer(ps, h.ID(), newHandle)
		} else {
			sendMessage(ps, msg)
		}

		// send msg with timer
		if strings.HasPrefix(msg, "/msgtimer ") {
			newTimer := strings.TrimPrefix(msg, "/msgtimer ")
			newTimerParam := strings.TrimSpace(newTimer)
			params := strings.Split(newTimerParam, " ")

			if len(params) < 1 {
				fmt.Println("Need two parameters: '50' '100' int number")
				return
			}

			// timer and loop
			timer, _ := strconv.Atoi(params[0])
			loop, _ := strconv.Atoi(params[1])

			// fake data with ammount chars
			msgMock := `transaction.Transaction{ID:"9e3b7881-f19f-4187-9290-6f5781c458c9", Sender:"rhzbc7082711414ab9823c7a5dedcb66edac4a2990f", Recipient:"rhz6a743d44af97c5b9b08311e78abea72e53e405e3", Signature:"20f57c312a759eeb52a86c18d927fb514e5d5234ff058c799539b24da103598c501369dd61b59f4771a0e0802f566a2e4a48939cd148daedde3e1c592b963780ba", Asset:map[string]string{}, Payload:[]uint8(nil), ContractData:[]uint8(nil), ContractResult:[]uint8(nil), Amount:amount.Amount{value:0xb2d05e00}, Fee:amount.Amount{value:0x0}, Type:0, VendorField:"", Timestamp:time.Time{wall:0x33f529a8, ext:63708734703, loc:(*time.Location)(nil)}, DoneData:transaction.DoneData{TxID:"", StartBlk:0x0, Timeout:0x0}}`

			for i := 0; i < loop; i++ {
				time.Sleep(time.Duration(timer) * time.Millisecond)
				sendMessage(ps, msgMock)
			}

		}

	}
	donec <- struct{}{}
}
