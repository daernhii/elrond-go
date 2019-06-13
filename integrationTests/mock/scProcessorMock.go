package mock

import (
	"github.com/ElrondNetwork/elrond-go-sandbox/data/state"
	"github.com/ElrondNetwork/elrond-go-sandbox/data/transaction"
	"github.com/ElrondNetwork/elrond-go-sandbox/process"
)

type SCProcessorMock struct {
	ComputeTransactionTypeCalled          func(tx *transaction.Transaction) (process.TransactionType, error)
	ExecuteSmartContractTransactionCalled func(tx *transaction.Transaction, acntSrc, acntDst state.AccountHandler, round uint32) error
	DeploySmartContractCalled             func(tx *transaction.Transaction, acntSrc state.AccountHandler, round uint32) error
}

func (sc *SCProcessorMock) ComputeTransactionType(
	tx *transaction.Transaction,
) (process.TransactionType, error) {
	if sc.ComputeTransactionTypeCalled == nil {
		return process.MoveBalance, nil
	}

	return sc.ComputeTransactionTypeCalled(tx)
}

func (sc *SCProcessorMock) ExecuteSmartContractTransaction(
	tx *transaction.Transaction,
	acntSrc, acntDst state.AccountHandler,
	round uint32,
) error {
	if sc.ExecuteSmartContractTransactionCalled == nil {
		return nil
	}

	return sc.ExecuteSmartContractTransactionCalled(tx, acntSrc, acntDst, round)
}

func (sc *SCProcessorMock) DeploySmartContract(
	tx *transaction.Transaction,
	acntSrc state.AccountHandler,
	round uint32,
) error {
	if sc.DeploySmartContractCalled == nil {
		return nil
	}

	return sc.DeploySmartContractCalled(tx, acntSrc, round)
}