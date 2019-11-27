package networksharding

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)

// NoSharder default sharder, only uses Kademlia distance in sorting
type NoSharder struct {
}

func (ns *NoSharder) ShouldReconnectToSeedNodes() bool {
	return false
}

func (ns *NoSharder) ShouldFindNewPeers() bool {
	return true
}

func (ns *NoSharder) CanConnectToPeer(current p2p.PeerID, connectedTo []p2p.PeerID, newPeer p2p.PeerID) bool {
	return true
}

func (ns *NoSharder) PeersForDisconnecting(current p2p.PeerID, connectedTo []p2p.PeerID) []p2p.PeerID {
	return make([]p2p.PeerID, 0)
}

// GetShard always 0
func (ns *NoSharder) GetShard(id peer.ID) uint32 {
	return 0
}

// GetDistance Kademlia XOR distance
func (ns *NoSharder) GetDistance(a, b sortingID) *big.Int {
	c := make([]byte, len(a.key))
	for i := 0; i < len(a.key); i++ {
		c[i] = a.key[i] ^ b.key[i]
	}

	ret := big.NewInt(0).SetBytes(c)
	return ret
}

// SortList sort the list
func (ns *NoSharder) SortList(peers []peer.ID, ref peer.ID) []peer.ID {
	return sortList(ns, peers, ref)
}

// IsInterfaceNil returns true if there is no value under the interface
func (ns *NoSharder) IsInterfaceNil() bool {
	return ns == nil
}
