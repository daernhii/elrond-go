package badgerdb_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/storage/badgerdb"
	"github.com/stretchr/testify/assert"
)

func createBadgerDb(t *testing.T, batchDelaySeconds int, maxBatchSize int) (p *badgerdb.DB) {
	dir, err := ioutil.TempDir("", "leveldb_temp")
	bDB, err := badgerdb.NewDB(dir, batchDelaySeconds, maxBatchSize)

	assert.Nil(t, err, "Failed creating leveldb database file")
	return bDB
}

func TestDB_InitNoError(t *testing.T) {
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Init()

	assert.Nil(t, err, "error initializing db")
}

func TestDB_PutNoError(t *testing.T) {
	key, val := []byte("key"), []byte("value")
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Put(key, val)

	assert.Nil(t, err, "error saving in db")
}

func TestDB_GetErrorAfterPutBeforeTimeout(t *testing.T) {
	key, val := []byte("key"), []byte("value")
	ldb := createBadgerDb(t, 1, 100)

	err := ldb.Put(key, val)
	assert.Nil(t, err)
	v, err := ldb.Get(key)
	assert.Nil(t, v)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Key not found")
}

func TestDB_GetOKAfterPutWithTimeout(t *testing.T) {
	key, val := []byte("key"), []byte("value")
	ldb := createBadgerDb(t, 1, 100)

	err := ldb.Put(key, val)
	assert.Nil(t, err)
	time.Sleep(time.Second * 3)

	v, err := ldb.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, val, v)
}

func TestDB_RemoveBeforeTimeoutOK(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	key, val := []byte("key"), []byte("value")
	ldb := createBadgerDb(t, 1, 100)

	err := ldb.Put(key, val)
	assert.Nil(t, err)

	_ = ldb.Remove(key)
	time.Sleep(time.Second * 2)

	v, err := ldb.Get(key)
	assert.Nil(t, v)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Key not found")
}

func TestDB_RemoveAfterTimeoutOK(t *testing.T) {
	key, val := []byte("key"), []byte("value")
	ldb := createBadgerDb(t, 1, 100)

	err := ldb.Put(key, val)
	assert.Nil(t, err)
	time.Sleep(time.Second * 2)

	_ = ldb.Remove(key)

	v, err := ldb.Get(key)
	assert.Nil(t, v)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Key not found")
}

func TestDB_GetPresent(t *testing.T) {
	key, val := []byte("key1"), []byte("value1")
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Put(key, val)

	assert.Nil(t, err, "error saving in db")

	v, err := ldb.Get(key)

	assert.Nil(t, err, "error not expected, but got %s", err)
	assert.Equalf(t, v, val, "read:%s but expected: %s", v, val)
}

func TestDB_GetNotPresent(t *testing.T) {
	key := []byte("key2")
	ldb := createBadgerDb(t, 10, 1)

	v, err := ldb.Get(key)

	assert.NotNil(t, err, "error expected but got nil, value %s", v)
}

func TestDB_HasPresent(t *testing.T) {
	key, val := []byte("key3"), []byte("value3")
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Put(key, val)

	assert.Nil(t, err, "error saving in db")

	err = ldb.Has(key)

	assert.Nil(t, err)
}

func TestDB_HasNotPresent(t *testing.T) {
	key := []byte("key4")
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Has(key)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Key not found")
}

func TestDB_RemovePresent(t *testing.T) {
	key, val := []byte("key5"), []byte("value5")
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Put(key, val)

	assert.Nil(t, err, "error saving in db")

	err = ldb.Remove(key)

	assert.Nil(t, err, "no error expected but got %s", err)

	err = ldb.Has(key)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Key not found")
}

func TestDB_RemoveNotPresent(t *testing.T) {
	key := []byte("key6")
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Remove(key)

	assert.Nil(t, err, "no error expected but got %s", err)
}

func TestDB_Close(t *testing.T) {
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Close()

	assert.Nil(t, err, "no error expected but got %s", err)
}

func TestDB_Destroy(t *testing.T) {
	ldb := createBadgerDb(t, 10, 1)

	err := ldb.Destroy()

	assert.Nil(t, err, "no error expected but got %s", err)
}
