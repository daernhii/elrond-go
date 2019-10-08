package wasm

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	"github.com/stretchr/testify/assert"
)

type mockAccountInfo struct {
	AddressBytes []byte
	Nonce        uint64
	Balance      *big.Int
}

func NewTestOwnerAccountInfo() mockAccountInfo {
	return mockAccountInfo{
		AddressBytes: []byte("12345678901234567890123456789012"),
		Nonce:        uint64(11),
		Balance:      big.NewInt(100000000),
	}
}

func TestWABTGasMetering_Off(t *testing.T) {
	fileSC := "./fibonacci_ewasmified.wasm"
	scOwnerAccount := NewTestOwnerAccountInfo()
	round := uint64(444)
	initialBalance := big.NewInt(1000)
	txValue := big.NewInt(3)

	aliceAccountInfo := mockAccountInfo{
		AddressBytes: []byte("12345678901234567890123456789111"),
		Nonce:        uint64(0),
		Balance:      initialBalance,
	}

	txProc, accnts := prepareTxProcessor(t, scOwnerAccount)
	_ = vm.CreateAccount(accnts, aliceAccountInfo.AddressBytes, aliceAccountInfo.Nonce, aliceAccountInfo.Balance)

	scAddress := deploySmartContract(t, txProc, accnts, fileSC, scOwnerAccount, round)
	callSC(t, scAddress, txProc, accnts, aliceAccountInfo, txValue, round, 0, 100)

	expectedBalance := big.NewInt(0).Sub(initialBalance, txValue)
	vm.TestAccount(t, accnts, aliceAccountInfo.AddressBytes, aliceAccountInfo.Nonce+1, expectedBalance)
}

func TestWABTGasMetering_On(t *testing.T) {
	fileSC := "./fibonacci_ewasmified.wasm"
	scOwnerAccount := NewTestOwnerAccountInfo()
	round := uint64(444)
	initialBalance := big.NewInt(1000)
	txValue := big.NewInt(3)

	aliceAccountInfo := mockAccountInfo{
		AddressBytes: []byte("12345678901234567890123456789111"),
		Nonce:        uint64(0),
		Balance:      initialBalance,
	}

	fmt.Printf("\n")
	fmt.Printf("Initial Balance of Alice: %s\n", aliceAccountInfo.Balance.String())
	fmt.Printf("\n")

	txProc, accnts := prepareTxProcessor(t, scOwnerAccount)
	_ = vm.CreateAccount(accnts, aliceAccountInfo.AddressBytes, aliceAccountInfo.Nonce, aliceAccountInfo.Balance)

	scAddress := deploySmartContract(t, txProc, accnts, fileSC, scOwnerAccount, round)
	callSC(t, scAddress, txProc, accnts, aliceAccountInfo, txValue, round, 1, 200)

	expectedSCExecutionCost := big.NewInt(76)

	expectedBalance := big.NewInt(0)
	expectedBalance.Set(initialBalance)
	expectedBalance.Sub(expectedBalance, txValue)
	expectedBalance.Sub(expectedBalance, expectedSCExecutionCost)

	vm.TestAccount(t, accnts, aliceAccountInfo.AddressBytes, aliceAccountInfo.Nonce+1, expectedBalance)
}

func callSC(
	t *testing.T,
	scAddress []byte,
	txProc process.TransactionProcessor,
	accnts state.AccountsAdapter,
	accountInfo mockAccountInfo,
	txValue *big.Int,
	round uint64,
	gasPrice uint64,
	gasLimit uint64,
) {
	tx := &transaction.Transaction{
		Nonce:     accountInfo.Nonce,
		Value:     txValue,
		RcvAddr:   scAddress,
		SndAddr:   accountInfo.AddressBytes,
		GasPrice:  gasPrice,
		GasLimit:  gasLimit,
		Data:      "benchmark",
		Signature: nil,
		Challenge: nil,
	}

	err := txProc.ProcessTransaction(tx, round)
	assert.Nil(t, err)
	_, err = accnts.Commit()
	assert.Nil(t, err)
}

func prepareTxProcessor(
	t *testing.T,
	initialAccount mockAccountInfo,
) (process.TransactionProcessor, state.AccountsAdapter) {
	txProc, accnts, _ := vm.CreatePreparedTxProcessorAndAccountsWithVMs(t, initialAccount.Nonce, initialAccount.AddressBytes, initialAccount.Balance)
	return txProc, accnts
}

func deploySmartContract(t *testing.T,
	txProc process.TransactionProcessor,
	accnts state.AccountsAdapter,
	fileSC string,
	owner mockAccountInfo,
	round uint64,
) []byte {
	scCode, err := ioutil.ReadFile(fileSC)
	assert.Nil(t, err)

	scCodeString := hex.EncodeToString(scCode)
	transferOnCalls := big.NewInt(1)
	gasPrice := uint64(1)
	gasLimit := uint64(100000)

	tx := vm.CreateTx(
		t,
		owner.AddressBytes,
		vm.CreateEmptyAddress().Bytes(),
		owner.Nonce,
		transferOnCalls,
		gasPrice,
		gasLimit,
		scCodeString+"@"+hex.EncodeToString(factory.HeraWABTVirtualMachine),
	)

	err = txProc.ProcessTransaction(tx, round)
	assert.Nil(t, err)

	_, err = accnts.Commit()
	assert.Nil(t, err)

	scAddress, _ := hex.DecodeString("000000000000000002001a2983b179a480a60c4308da48f13b4480dbb4d33132")
	return scAddress
}
