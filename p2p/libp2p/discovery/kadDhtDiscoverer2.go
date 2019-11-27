package discovery

import (
	"context"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

// KadDhtDiscoverer2 is the kad-dht discovery type implementation
type KadDhtDiscoverer2 struct {
	mutBootstrapStarted sync.Mutex
	bootstrapStarted    bool
	kadDHT              ipfsDiscoverer
	connHost            libp2p.ConnectableHost
	ctx                 context.Context
	sharder             p2p.Sharder

	refreshInterval time.Duration
	randezVous      string
	seeders         []string
}

// NewKadDhtPeerDiscoverer2 creates a new kad-dht discovery type implementation
// initialPeersList can be nil or empty, no initial connection will be attempted, a warning message will appear
func NewKadDhtPeerDiscoverer2(
	kadDHT ipfsDiscoverer,
	connHost libp2p.ConnectableHost,
	ctx context.Context,
	sharder p2p.Sharder,
	refreshInterval time.Duration,
	randezVous string,
	initialPeersList []string,
) (*KadDhtDiscoverer2, error) {

	if kadDHT == nil {
		return nil, p2p.ErrNilPeerDiscoverer
	}
	if check.IfNil(connHost) {
		return nil, p2p.ErrNilHost
	}
	if ctx == nil {
		return nil, p2p.ErrNilContext
	}
	if check.IfNil(sharder) {
		return nil, p2p.ErrNilSharder
	}

	if len(initialPeersList) == 0 {
		log.Warn("nil or empty initial peers list provided to kad dht implementation. " +
			"No initial connection will be done")
	}

	return &KadDhtDiscoverer2{
		kadDHT:          kadDHT,
		connHost:        connHost,
		ctx:             ctx,
		sharder:         sharder,
		refreshInterval: refreshInterval,
		randezVous:      randezVous,
		seeders:         initialPeersList,
	}, nil
}

// Bootstrap will start the bootstrapping new peers process
func (kdd *KadDhtDiscoverer2) Bootstrap() error {
	kdd.mutBootstrapStarted.Lock()
	defer kdd.mutBootstrapStarted.Unlock()

	if kdd.bootstrapStarted {
		return p2p.ErrPeerDiscoveryProcessAlreadyStarted
	}

	go kdd.doBootstrap()

	return nil
}

func (kdd *KadDhtDiscoverer2) doBootstrap() {
	for {
		if kdd.shouldReconnectToSeednodes() {
			kdd.reconectToSeedNodes()
		}

		if kdd.sharder.ShouldFindNewPeers() {
			kdd.findPeers()
		}

		select {
		case <-time.After(kdd.refreshInterval):
		case <-kdd.ctx.Done():
			return
		}
	}
}

func (kdd *KadDhtDiscoverer2) shouldReconnectToSeednodes() bool {
	if len(kdd.seeders) == 0 {
		return false
	}

	connectedSeeders := 0
	for _, seeder := range kdd.seeders {
		if kdd.isConnectedToPeer(seeder) {
			connectedSeeders++
		}
	}

	if kdd.sharder.ShouldReconnectToSeedNodes() {
		return true
	}

	return connectedSeeders < len(kdd.seeders)
}

func (kdd *KadDhtDiscoverer2) isConnectedToPeer(address string) bool {
	multiAddr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return false
	}

	pInfo, err := peer.AddrInfoFromP2pAddr(multiAddr)
	if err != nil {
		return false
	}

	connectedness := kdd.connHost.Network().Connectedness(pInfo.ID)

	return connectedness == network.Connected
}

func (kdd *KadDhtDiscoverer2) reconectToSeedNodes() {
	if len(kdd.seeders) == 0 {
		return
	}

	for _, seeder := range kdd.seeders {
		if !kdd.isConnectedToPeer(seeder) {
			err := kdd.connHost.ConnectToPeer(kdd.ctx, seeder)
			if err != nil {
				log.Debug("error connecting to seeder", "error", err.Error())
				continue
			}
		}
	}
}

func (kdd *KadDhtDiscoverer2) findPeers() {
	cfg := dht.BootstrapConfig{
		Period:  kdd.refreshInterval,
		Queries: noOfQueries,
		Timeout: peerDiscoveryTimeout,
	}

	err := kdd.kadDHT.BootstrapOnce(kdd.ctx, cfg)
	if err != nil {
		log.Debug("error finding peers", "error", err.Error())
	}
}

// Name returns the name of the kad dht peer discovery implementation
func (kdd *KadDhtDiscoverer2) Name() string {
	return kadDhtName
}

// ApplyContext is DEPRECATED
func (kdd *KadDhtDiscoverer2) ApplyContext(ctxProvider p2p.ContextProvider) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (kdd *KadDhtDiscoverer2) IsInterfaceNil() bool {
	return kdd == nil
}
