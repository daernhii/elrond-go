package mock

import "github.com/ElrondNetwork/elrond-go/p2p"

type SharderStub struct {
	ShouldReconnectToSeedNodesCalled func() bool
	ShouldFindNewPeersCalled         func() bool
	CanConnectToPeerCalled           func(current p2p.PeerID, connectedTo []p2p.PeerID, newPeer p2p.PeerID) bool
	PeersForDisconnectingCalled      func(current p2p.PeerID, connectedTo []p2p.PeerID) []p2p.PeerID
}

func (ss *SharderStub) ShouldReconnectToSeedNodes() bool {
	return ss.ShouldReconnectToSeedNodesCalled()
}

func (ss *SharderStub) ShouldFindNewPeers() bool {
	return ss.ShouldFindNewPeersCalled()
}

func (ss *SharderStub) CanConnectToPeer(current p2p.PeerID, connectedTo []p2p.PeerID, newPeer p2p.PeerID) bool {
	return ss.CanConnectToPeerCalled(current, connectedTo, newPeer)
}

func (ss *SharderStub) PeersForDisconnecting(current p2p.PeerID, connectedTo []p2p.PeerID) []p2p.PeerID {
	return ss.PeersForDisconnectingCalled(current, connectedTo)
}

func (ss *SharderStub) IsInterfaceNil() bool {
	return ss == nil
}
