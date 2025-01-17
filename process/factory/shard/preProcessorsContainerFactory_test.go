package shard

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewPreProcessorsContainerFactory_NilShardCoordinator(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		nil,
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilShardCoordinator, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilStore(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		nil,
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilStore, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilMarshalizer(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		nil,
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilMarshalizer, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilHasher(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		nil,
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilHasher, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilDataPool(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		nil,
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilDataPoolHolder, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilAddrConv(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		nil,
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilAddressConverter, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilAccounts(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		nil,
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilAccountsAdapter, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilTxProcessor(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		nil,
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilTxProcessor, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilSCProcessor(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		nil,
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilSmartContractProcessor, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilSCR(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		nil,
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilSmartContractResultProcessor, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilRewardTxProcessor(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		nil,
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilRewardsTxProcessor, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilRequestHandler(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		nil,
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilRequestHandler, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilInternalTransactionProducer(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		nil,
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilInternalTransactionProducer, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilFeeHandler(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		nil,
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilEconomicsFeeHandler, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilMiniBlocksCompacter(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		nil,
		&mock.GasHandlerMock{},
	)

	assert.Equal(t, process.ErrNilMiniBlocksCompacter, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory_NilGasHandler(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		nil,
	)

	assert.Equal(t, process.ErrNilGasHandler, err)
	assert.Nil(t, ppcm)
}

func TestNewPreProcessorsContainerFactory(t *testing.T) {
	t.Parallel()

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		mock.NewPoolsHolderMock(),
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, ppcm)
}

func TestPreProcessorsContainerFactory_CreateErrTxPreproc(t *testing.T) {
	t.Parallel()
	dataPool := &mock.PoolsHolderStub{}
	dataPool.TransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return nil
	}

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		dataPool,
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, ppcm)

	container, err := ppcm.Create()
	assert.Nil(t, container)
	assert.Equal(t, process.ErrNilTransactionPool, err)
}

func TestPreProcessorsContainerFactory_CreateErrScrPreproc(t *testing.T) {
	t.Parallel()
	dataPool := &mock.PoolsHolderStub{}
	dataPool.TransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return &mock.ShardedDataStub{
			RegisterHandlerCalled: func(i func(key []byte)) {
			},
		}
	}
	dataPool.UnsignedTransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return nil
	}

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		dataPool,
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, ppcm)

	container, err := ppcm.Create()
	assert.Nil(t, container)
	assert.Equal(t, process.ErrNilUTxDataPool, err)
}

func TestPreProcessorsContainerFactory_Create(t *testing.T) {
	t.Parallel()
	dataPool := &mock.PoolsHolderStub{}
	dataPool.TransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return &mock.ShardedDataStub{
			RegisterHandlerCalled: func(i func(key []byte)) {
			},
		}
	}
	dataPool.UnsignedTransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return &mock.ShardedDataStub{
			RegisterHandlerCalled: func(i func(key []byte)) {
			},
		}
	}
	dataPool.RewardTransactionsCalled = func() dataRetriever.ShardedDataCacherNotifier {
		return &mock.ShardedDataStub{
			RegisterHandlerCalled: func(i func(key []byte)) {
			},
		}
	}

	ppcm, err := NewPreProcessorsContainerFactory(
		mock.NewMultiShardsCoordinatorMock(3),
		&mock.ChainStorerMock{},
		&mock.MarshalizerMock{},
		&mock.HasherMock{},
		dataPool,
		&mock.AddressConverterMock{},
		&mock.AccountsStub{},
		&mock.RequestHandlerMock{},
		&mock.TxProcessorMock{},
		&mock.SCProcessorMock{},
		&mock.SmartContractResultsProcessorMock{},
		&mock.RewardTxProcessorMock{},
		&mock.IntermediateTransactionHandlerMock{},
		&mock.FeeHandlerStub{},
		&mock.MiniBlocksCompacterMock{},
		&mock.GasHandlerMock{},
	)

	assert.Nil(t, err)
	assert.NotNil(t, ppcm)

	container, err := ppcm.Create()
	assert.Nil(t, err)
	assert.Equal(t, 3, container.Len())
}
