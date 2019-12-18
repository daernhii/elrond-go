package smartContractResult_test

import (
	"bytes"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/smartContractResult"
	"github.com/stretchr/testify/assert"
)

func TestSmartContractResult_SaveLoad(t *testing.T) {
	smrS := smartContractResult.SmartContractResult{
		Nonce:   uint64(1),
		Value:   data.NewProtoBigInt(1),
		RcvAddr: []byte("receiver_address"),
		SndAddr: []byte("sender_address"),
		Data:    "scr_data",
		Code:    []byte("code"),
		TxHash:  []byte("scrHash"),
	}

	var b bytes.Buffer
	_ = smrS.Save(&b)

	loadSMR := smartContractResult.SmartContractResult{}
	_ = loadSMR.Load(&b)

	assert.Equal(t, smrS, loadSMR)
}

func TestSmartContractResult_GetData(t *testing.T) {
	t.Parallel()

	data := "data"
	scr := &smartContractResult.SmartContractResult{Data: data}

	assert.Equal(t, data, scr.Data)
}

func TestSmartContractResult_GetRecvAddr(t *testing.T) {
	t.Parallel()

	data := []byte("data")
	scr := &smartContractResult.SmartContractResult{RcvAddr: data}

	assert.Equal(t, data, scr.RcvAddr)
}

func TestSmartContractResult_GetSndAddr(t *testing.T) {
	t.Parallel()

	data := []byte("data")
	scr := &smartContractResult.SmartContractResult{SndAddr: data}

	assert.Equal(t, data, scr.SndAddr)
}

func TestSmartContractResult_GetValue(t *testing.T) {
	t.Parallel()

	value := data.NewProtoBigInt(10)
	scr := &smartContractResult.SmartContractResult{Value: value}

	assert.Equal(t, value, scr.Value)
}

func TestSmartContractResult_SetData(t *testing.T) {
	t.Parallel()

	data := "data"
	scr := &smartContractResult.SmartContractResult{}
	scr.SetData(data)

	assert.Equal(t, data, scr.Data)
}

func TestSmartContractResult_SetRecvAddr(t *testing.T) {
	t.Parallel()

	data := []byte("data")
	scr := &smartContractResult.SmartContractResult{}
	scr.SetRcvAddr(data)

	assert.Equal(t, data, scr.RcvAddr)
}

func TestSmartContractResult_SetSndAddr(t *testing.T) {
	t.Parallel()

	data := []byte("data")
	scr := &smartContractResult.SmartContractResult{}
	scr.SetSndAddr(data)

	assert.Equal(t, data, scr.SndAddr)
}

func TestSmartContractResult_SetValue(t *testing.T) {
	t.Parallel()

	value := data.NewProtoBigInt(10)
	scr := &smartContractResult.SmartContractResult{}
	scr.SetValue(value)

	assert.Equal(t, value, scr.Value)
}
