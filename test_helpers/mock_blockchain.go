// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test_helpers

import (
	"math/big"

	"github.com/vulcanize/vulcanizedb/pkg/core"
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

func (MockBlockchain) GetHeaderByNumber(blockNumber int64) (core.Header, error) {
	panic("implement me")
}

func (MockBlockchain) GetBlockByNumber(blockNumber int64) (core.Block, error) {
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
