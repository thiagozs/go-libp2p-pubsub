package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"strings"

	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// ChatInputLoop handler for message between peers
func ChatInputLoop(cc *counters, ctx context.Context, h host.Host, ps *pubsub.PubSub, donec chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()

		switch {
		case reset(cc, msg):
		case showStats(cc, msg):
		case changeName(msg, h, ps):
		case messageWithTimer(msg, ps):
		default:
			if len(strings.TrimSpace(msg)) > 0 {
				sendMessage(ps, msg)
			}
		}
	}
	donec <- struct{}{}
}

func reset(cc *counters, msg string) bool {
	if strings.HasPrefix(msg, "/reset") {
		cc.Reset()
		return true
	}
	return false
}

func showStats(cc *counters, msg string) bool {
	if strings.HasPrefix(msg, "/stats") {
		cc.Show()
		return true
	}
	return false
}

func changeName(msg string, h host.Host, ps *pubsub.PubSub) bool {
	if strings.HasPrefix(msg, "/name") {
		newHandle := strings.TrimPrefix(msg, "/name ")
		newHandle = strings.TrimSpace(newHandle)
		updatePeer(ps, h.ID(), newHandle)
		return true
	}
	return false
}

func messageWithTimer(msg string, ps *pubsub.PubSub) bool {
	if strings.HasPrefix(msg, "/msgtimer") {
		newTimer := strings.TrimPrefix(msg, "/msgtimer ")
		newTimerParam := strings.TrimSpace(newTimer)
		params := strings.Split(newTimerParam, " ")

		// check parameters
		if len(params) < 2 {
			fmt.Println("Need two parameters: '50' '100' 'slug' int number")
		} else {
			// timer and loop
			timer, _ := strconv.Atoi(params[0])
			loop, _ := strconv.Atoi(params[1])

			// fake data with ammount chars
			msgMock := `transaction.Transaction{ID:"9e3b7881-f19f-4187-9290-6f5781c458c9", Sender:"rhzbc7082711414ab9823c7a5dedcb66edac4a2990f", Recipient:"rhz6a743d44af97c5b9b08311e78abea72e53e405e3", Signature:"20f57c312a759eeb52a86c18d927fb514e5d5234ff058c799539b24da103598c501369dd61b59f4771a0e0802f566a2e4a48939cd148daedde3e1c592b963780ba", Asset:map[string]string{}, Payload:[]uint8(nil), ContractData:[]uint8(nil), ContractResult:[]uint8(nil), Amount:amount.Amount{value:0xb2d05e00}, Fee:amount.Amount{value:0x0}, Type:0, VendorField:"", Timestamp:time.Time{wall:0x33f529a8, ext:63708734703, loc:(*time.Location)(nil)}, DoneData:transaction.DoneData{TxID:"", StartBlk:0x0, Timeout:0x0}}`

			// send without lock terminal
			go func() {
				for i := 0; i < loop; i++ {
					time.Sleep(time.Duration(timer) * time.Millisecond)
					sendMessage(ps, fmt.Sprintf("%s| %s - %s", params[2], strconv.Itoa(i), msgMock))
				}
			}()
		}
		return true
	}
	return false
}
