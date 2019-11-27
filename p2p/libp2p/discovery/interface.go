package discovery

import (
	"context"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type ipfsDiscoverer interface {
	BootstrapOnce(ctx context.Context, config dht.BootstrapConfig) error
}
