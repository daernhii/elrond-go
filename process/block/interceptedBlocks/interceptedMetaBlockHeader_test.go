package interceptedBlocks_test

import (
	"bytes"
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/crypto"
	dataBlock "github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/block/interceptedBlocks"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/stretchr/testify/assert"
)

func createDefaultMetaArgument() *interceptedBlocks.ArgInterceptedBlockHeader {
	arg := &interceptedBlocks.ArgInterceptedBlockHeader{
		ShardCoordinator: mock.NewOneShardCoordinatorMock(),
		MultiSigVerifier: mock.NewMultiSigner(),
		Hasher:           testHasher,
		Marshalizer:      testMarshalizer,
		NodesCoordinator: &mock.NodesCoordinatorMock{
			ComputeValidatorsGroupCalled: func(randomness []byte, round uint64, shardId uint32) (validatorsGroup []sharding.Validator, err error) {

				validator := mock.NewValidatorMock(big.NewInt(0), 0, []byte("pubKey"), []byte("pubKey"))
				return []sharding.Validator{validator}, nil
			},
		},
		KeyGen: &mock.SingleSignKeyGenMock{
			PublicKeyFromByteArrayCalled: func(b []byte) (key crypto.PublicKey, err error) {
				return nil, nil
			},
		},
		SingleSigVerifier: &mock.SignerMock{
			VerifyStub: func(public crypto.PublicKey, msg []byte, sig []byte) error {
				return nil
			},
		},
	}

	hdr := createMockMetaHeader()
	arg.HdrBuff, _ = testMarshalizer.Marshal(hdr)

	return arg
}

func createMockMetaHeader() *dataBlock.MetaBlock {
	return &dataBlock.MetaBlock{
		Nonce:         hdrNonce,
		PrevHash:      []byte("prev hash"),
		PrevRandSeed:  []byte("prev rand seed"),
		RandSeed:      []byte("rand seed"),
		PubKeysBitmap: []byte{1},
		TimeStamp:     0,
		Round:         hdrRound,
		Epoch:         hdrEpoch,
		Signature:     []byte("signature"),
		RootHash:      []byte("root hash"),
		TxCount:       0,
		PeerInfo:      nil,
		ShardInfo:     nil,
	}
}

//------- TestNewInterceptedHeader

func TestNewInterceptedMetaHeader_NilArgumentShouldErr(t *testing.T) {
	t.Parallel()

	inHdr, err := interceptedBlocks.NewInterceptedMetaHeader(nil)

	assert.Nil(t, inHdr)
	assert.Equal(t, process.ErrNilArguments, err)
}

func TestNewInterceptedMetaHeader_MarshalizerFailShouldErr(t *testing.T) {
	t.Parallel()

	arg := createDefaultMetaArgument()
	arg.HdrBuff = []byte("invalid buffer")

	inHdr, err := interceptedBlocks.NewInterceptedMetaHeader(arg)

	assert.Nil(t, inHdr)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestNewInterceptedMetaHeader_ShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultMetaArgument()

	inHdr, err := interceptedBlocks.NewInterceptedMetaHeader(arg)

	assert.False(t, check.IfNil(inHdr))
	assert.Nil(t, err)
}

//------- CheckValidity

func TestInterceptedMetaHeader_CheckValidityNilPubKeyBitmapShouldErr(t *testing.T) {
	t.Parallel()

	hdr := createMockMetaHeader()
	hdr.PubKeysBitmap = nil
	buff, _ := testMarshalizer.Marshal(hdr)

	arg := createDefaultMetaArgument()
	arg.HdrBuff = buff
	inHdr, _ := interceptedBlocks.NewInterceptedMetaHeader(arg)

	err := inHdr.CheckValidity()

	assert.Equal(t, process.ErrNilPubKeysBitmap, err)
}

func TestInterceptedMetaHeader_ErrorInMiniBlockShouldErr(t *testing.T) {
	t.Parallel()

	hdr := createMockMetaHeader()
	badShardId := uint32(2)
	hdr.ShardInfo = []dataBlock.ShardData{
		{
			ShardID:               badShardId,
			HeaderHash:            nil,
			ShardMiniBlockHeaders: nil,
			TxCount:               0,
		},
	}
	buff, _ := testMarshalizer.Marshal(hdr)

	arg := createDefaultShardArgument()
	arg.HdrBuff = buff
	inHdr, _ := interceptedBlocks.NewInterceptedMetaHeader(arg)

	err := inHdr.CheckValidity()

	assert.Equal(t, process.ErrInvalidShardId, err)
}

func TestInterceptedMetaHeader_CheckValidityShouldWork(t *testing.T) {
	t.Parallel()

	arg := createDefaultMetaArgument()
	inHdr, _ := interceptedBlocks.NewInterceptedMetaHeader(arg)

	err := inHdr.CheckValidity()

	assert.Nil(t, err)
}

//------- getters

func TestInterceptedMetaHeader_Getters(t *testing.T) {
	t.Parallel()

	arg := createDefaultMetaArgument()
	inHdr, _ := interceptedBlocks.NewInterceptedMetaHeader(arg)

	hash := testHasher.Compute(string(arg.HdrBuff))

	assert.Equal(t, hash, inHdr.Hash())
	assert.True(t, inHdr.IsForCurrentShard())
}

func TestInterceptedMetaHeader_CheckValidityLeaderSignatureNotCorrectShouldErr(t *testing.T) {
	t.Parallel()

	hdr := createMockShardHeader()
	expectedErr := errors.New("expected err")
	buff, _ := testMarshalizer.Marshal(hdr)

	arg := createDefaultShardArgument()
	arg.SingleSigVerifier = &mock.SignerMock{
		SignStub: nil,
		VerifyStub: func(public crypto.PublicKey, msg []byte, sig []byte) error {
			return expectedErr
		},
	}
	arg.HdrBuff = buff
	inHdr, _ := interceptedBlocks.NewInterceptedMetaHeader(arg)

	err := inHdr.CheckValidity()
	assert.Equal(t, expectedErr, err)
}

func TestInterceptedMetaHeader_CheckValidityLeaderSignatureOkShouldWork(t *testing.T) {
	t.Parallel()

	hdr := createMockShardHeader()
	expectedSignature := []byte("ran")
	hdr.LeaderSignature = expectedSignature
	buff, _ := testMarshalizer.Marshal(hdr)

	arg := createDefaultShardArgument()
	arg.SingleSigVerifier = &mock.SignerMock{
		SignStub: nil,
		VerifyStub: func(public crypto.PublicKey, msg []byte, sig []byte) error {
			// skip this for signature check. only leader's signature is relevant for this test
			if !bytes.Equal(sig, []byte("rand seed")) {
				isSignOk := bytes.Equal(sig, expectedSignature)
				assert.True(t, isSignOk)
			}
			return nil
		},
	}
	arg.HdrBuff = buff
	inHdr, _ := interceptedBlocks.NewInterceptedMetaHeader(arg)

	err := inHdr.CheckValidity()
	assert.Nil(t, err)
}

//------- IsInterfaceNil

func TestInterceptedMetaHeader_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var inHdr *interceptedBlocks.InterceptedMetaHeader

	assert.True(t, check.IfNil(inHdr))
}
