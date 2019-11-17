package networksharding

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/libp2p/go-libp2p-core/peer"
)

type lastByteResolver struct{}

func (res *lastByteResolver) ByID(pid p2p.PeerID) uint32 {
	if pid[len(pid)-1] == 'M' {
		return sharding.MetachainShardId
	}

	if pid[len(pid)-1] == 'U' {
		return unknownShardId
	}

	return uint32(pid[len(pid)-1] - 48)
}

func (res *lastByteResolver) IsInterfaceNil() bool {
	return res == nil
}

func transformStringsToPeerIds(peers []string) []peer.ID {
	result := make([]peer.ID, len(peers))
	for i, p := range peers {
		result[i] = peer.ID(p)
	}

	return result
}

func TestKadSharderWithLists_GetShard(t *testing.T) {
	t.Parallel()

	kswl := NewKadSharderWithLists(&lastByteResolver{})

	pids := []string{"keyShard0", "keyShard1", "keyShard2", "keyShard9", "keyShardM", "keyShardU"}

	for i := 0; i < len(pids); i++ {
		fmt.Printf("%s : shard %d\n", pids[i], kswl.GetShard(peer.ID(pids[i])))
	}
}

func TestKadSharderWithLists_SortList(t *testing.T) {
	t.Parallel()

	kswl := NewKadSharderWithLists(&lastByteResolver{})
	peers := []string{"keyShard0", "keyShard1", "anotherKeyShard0", "key1Shard1", "peerShard0", "keyShard2", "keyShard9", "keyShardM", "keyShardU", "another1KeyShard0"}

	result := kswl.SortList(transformStringsToPeerIds(peers), peer.ID("peerShard0"))

	for _, r := range result {
		fmt.Printf("%s\n", string(r))
	}

	fmt.Println()

	result2 := kswl.SortList(transformStringsToPeerIds(peers), peer.ID("anotherKeyShard0"))

	for _, r := range result2 {
		fmt.Printf("%s\n", string(r))
	}
}
