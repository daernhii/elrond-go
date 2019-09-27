package whiteList

import "github.com/ElrondNetwork/elrond-go/storage/lrucache"

//WhiteListData will hold data from other shard
//will hold transactions hashes and mini blocks hashes
type WhiteListData struct {
	dataHashes *lrucache.LRUCache
}

// NewWhiteListData will return an object that stores a white list
func NewWhiteListData(sizeCache int) (*WhiteListData, error) {

	dataHashes, err := lrucache.NewCache(sizeCache)
	if err != nil {
		return nil, err
	}

	return &WhiteListData{
		dataHashes: dataHashes,
	}, nil
}

// AddHash will add a hash in white list
func (wld *WhiteListData) AddHash(hash []byte) {
	wld.dataHashes.Put(hash, struct{}{})
}

// RemoveHash will remove a hash in white list
func (wld *WhiteListData) RemoveHash(hash []byte) {
	wld.dataHashes.Remove(hash)
}

// RemoveHashes will remove a list of hashed from white list
func (wld *WhiteListData) RemoveHashes(hashes [][]byte) {
	for _, hash := range hashes {
		wld.dataHashes.Remove(hash)
	}
}

// IsInWhiteList will check if in white list exits a hash
func (wld *WhiteListData) IsInWhiteList(hash []byte) bool {
	return wld.dataHashes.Has(hash)
}
