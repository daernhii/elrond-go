package mock

import (
	"github.com/ElrondNetwork/elrond-go/consensus"
	"github.com/ElrondNetwork/elrond-go/data"
)

type BroadcastMessengerMock struct {
	BroadcastBlockCalled            func(data.BodyHandler, data.HeaderHandler) error
	BroadcastHeaderCalled           func(data.HeaderHandler) error
	BroadcastShardHeaderCalled      func(data.HeaderHandler) error
	BroadcastMiniBlocksCalled       func(map[uint32][]byte) error
	BroadcastTransactionsCalled     func(map[string][][]byte) error
	BroadcastConsensusMessageCalled func(*consensus.Message) error
}

func (bmm *BroadcastMessengerMock) BroadcastBlock(bodyHandler data.BodyHandler, headerhandler data.HeaderHandler) error {
	if bmm.BroadcastBlockCalled != nil {
		return bmm.BroadcastBlockCalled(bodyHandler, headerhandler)
	}
	return nil
}

func (bmm *BroadcastMessengerMock) BroadcastShardHeader(headerHandler data.HeaderHandler) error {
	if bmm.BroadcastShardHeaderCalled != nil {
		return bmm.BroadcastShardHeaderCalled(headerHandler)
	}
	return nil
}

func (bmm *BroadcastMessengerMock) BroadcastMiniBlocks(miniBlocks map[uint32][]byte) error {
	if bmm.BroadcastMiniBlocksCalled != nil {
		return bmm.BroadcastMiniBlocksCalled(miniBlocks)
	}
	return nil
}

func (bmm *BroadcastMessengerMock) BroadcastTransactions(transactions map[string][][]byte) error {
	if bmm.BroadcastTransactionsCalled != nil {
		return bmm.BroadcastTransactionsCalled(transactions)
	}
	return nil
}

func (bmm *BroadcastMessengerMock) BroadcastConsensusMessage(message *consensus.Message) error {
	if bmm.BroadcastConsensusMessageCalled != nil {
		return bmm.BroadcastConsensusMessageCalled(message)
	}
	return nil
}

func (bmm *BroadcastMessengerMock) BroadcastHeader(headerhandler data.HeaderHandler) error {
	if bmm.BroadcastHeaderCalled != nil {
		return bmm.BroadcastHeaderCalled(headerhandler)
	}
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (bmm *BroadcastMessengerMock) IsInterfaceNil() bool {
	if bmm == nil {
		return true
	}
	return false
}
