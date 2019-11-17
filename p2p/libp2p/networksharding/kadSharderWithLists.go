package networksharding

import (
	"math/big"
	"sort"

	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/libp2p/go-libp2p-core/peer"
)

const unknownShardId = uint32(0xFFFFFFFE)

type kadSharderWithLists struct {
	resolver p2p.PeerShardResolver
}

func NewKadSharderWithLists(resolver p2p.PeerShardResolver) *kadSharderWithLists {
	return &kadSharderWithLists{
		resolver: resolver,
	}
}

func (kswl *kadSharderWithLists) GetShard(id peer.ID) uint32 {
	return kswl.resolver.ByID(p2p.PeerID(id))
}

func (kswl *kadSharderWithLists) GetDistance(a, b sortingID) *big.Int {
	c := make([]byte, len(a.key))
	for i := 0; i < len(a.key); i++ {
		c[i] = a.key[i] ^ b.key[i]
	}

	ret := big.NewInt(0).SetBytes(c)
	return ret
}

func (kswl *kadSharderWithLists) SortList(peers []peer.ID, ref peer.ID) []peer.ID {
	currentSid := sortingID{
		id:       ref,
		key:      keyFromID(ref),
		shard:    kswl.GetShard(ref),
		distance: big.NewInt(0),
	}

	sameShardList, crossShardList := kswl.splitPeers(peers, currentSid)
	kswl.sortList(sameShardList)
	kswl.sortList(crossShardList)

	return kswl.mergeLists(sameShardList, crossShardList)
}

func (kswl *kadSharderWithLists) splitPeers(peers []peer.ID, currentSid sortingID) ([]sortingID, []sortingID) {
	sameShard := make([]sortingID, 0)
	crossShard := make([]sortingID, 0)

	for _, p := range peers {
		peerShard := kswl.GetShard(p)
		isCrossShard := peerShard != currentSid.shard || peerShard == unknownShardId

		sid := sortingID{
			id:    p,
			key:   keyFromID(p),
			shard: kswl.GetShard(p),
		}
		sid.distance = kswl.GetDistance(sid, currentSid)

		if isCrossShard {
			crossShard = append(crossShard, sid)
		} else {
			sameShard = append(sameShard, sid)
		}
	}

	return sameShard, crossShard
}

func (kswl *kadSharderWithLists) sortList(list []sortingID) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].distance.Cmp(list[j].distance) < 0
	})
}

func (kswl *kadSharderWithLists) mergeLists(sameShardList []sortingID, crossShardList []sortingID) []peer.ID {
	result := make([]peer.ID, 0)

	i := 0
	for ; i < len(sameShardList) && i < len(crossShardList); i++ {
		result = append(result, sameShardList[i].id)
		result = append(result, crossShardList[i].id)
	}

	for ; i < len(sameShardList); i++ {
		result = append(result, sameShardList[i].id)
	}

	for ; i < len(crossShardList); i++ {
		result = append(result, crossShardList[i].id)
	}

	return result
}

func (kswl *kadSharderWithLists) IsInterfaceNil() bool {
	return kswl == nil
}
