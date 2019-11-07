package mdns

import (
	"context"

	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Notifee dns system notifyer
type Notifee struct {
	H   host.Host
	CTX context.Context
}

// HandlePeerFound find peer handker
func (m *Notifee) HandlePeerFound(pi peer.AddrInfo) {
	_ = m.H.Connect(m.CTX, pi)
}
