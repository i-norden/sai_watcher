package test_helpers

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/big"
)

type MockBlockchain struct {
	AbiJSONs     []string
	Addresses    []string
	Methods      []string
	MethodArgs   []interface{}
	Results      []interface{}
	BlockNumbers []int64
}

func (mb *MockBlockchain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	mb.AbiJSONs = append(mb.AbiJSONs, abiJSON)
	mb.Addresses = append(mb.Addresses, address)
	mb.Methods = append(mb.Methods, method)
	mb.MethodArgs = append(mb.MethodArgs, methodArg)
	mb.Results = append(mb.Results, result)
	mb.BlockNumbers = append(mb.BlockNumbers, blockNumber)
	return nil
}

func (MockBlockchain) GetBlockByNumber(blockNumber int64) core.Block {
	panic("implement me")
}

func (MockBlockchain) GetLogs(contract core.Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]core.Log, error) {
	panic("implement me")
}

func (MockBlockchain) LastBlock() *big.Int {
	panic("implement me")
}

func (MockBlockchain) Node() core.Node {
	panic("implement me")
}
