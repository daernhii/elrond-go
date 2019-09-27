package whiteList_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/process/whiteList"
	"github.com/stretchr/testify/assert"
)

func TestWhiteListData_ZeroCacheSizeShouldErr(t *testing.T) {
	t.Parallel()

	whiteListData, err := whiteList.NewWhiteListData(0)

	assert.Equal(t, errors.New("Must provide a positive size"), err)
	assert.Nil(t, whiteListData)
}

func TestWhiteListData_ShouldWork(t *testing.T) {
	t.Parallel()

	whiteListData, err := whiteList.NewWhiteListData(1)

	assert.Nil(t, err)
	assert.NotNil(t, whiteListData)
}

func TestWhiteListData_AddHash(t *testing.T) {
	t.Parallel()

	hash := []byte("hash")
	whiteListData, _ := whiteList.NewWhiteListData(1)

	whiteListData.AddHash(hash)
	hashWasAdded := whiteListData.IsInWhiteList(hash)

	assert.Equal(t, true, hashWasAdded)
}

func TestWhiteListData_RemoveHash(t *testing.T) {
	t.Parallel()

	hash := []byte("hash")
	whiteListData, _ := whiteList.NewWhiteListData(1)

	whiteListData.AddHash(hash)
	whiteListData.RemoveHash(hash)
	hashWasAdded := whiteListData.IsInWhiteList(hash)

	assert.Equal(t, false, hashWasAdded)
}

func TestWhiteListData_RemoveHashes(t *testing.T) {
	t.Parallel()

	hash1 := []byte("hash1")
	hash2 := []byte("hash2")
	whiteListData, _ := whiteList.NewWhiteListData(2)

	whiteListData.AddHash(hash1)
	whiteListData.AddHash(hash2)
	whiteListData.RemoveHashes([][]byte{hash1, hash2})

	hashWasAdded := whiteListData.IsInWhiteList(hash1)
	assert.Equal(t, false, hashWasAdded)
	hashWasAdded = whiteListData.IsInWhiteList(hash2)
	assert.Equal(t, false, hashWasAdded)
}
