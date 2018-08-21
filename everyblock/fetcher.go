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

package everyblock

import (
	"math/big"

	"github.com/8thlight/sai_watcher/event_triggered/tub"
	"github.com/8thlight/sai_watcher/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

func NewFetcher(blockchain core.BlockChain) *Fetcher {
	return &Fetcher{
		blockchain: blockchain,
	}
}

type Fetcher struct {
	blockchain core.BlockChain
}

type Value [32]byte

type Peek struct {
	Value
	OK bool
}

func newResult(value [32]byte, ok bool) *Peek {
	return &Peek{Value: value, OK: ok}
}

func (peek Peek) Wad() string {
	return utils.Convert("wad", peek.Value.String(), 15)
}

func (value Value) Hex() string {
	return hexutil.Encode(value[:])
}

func (value Value) String() string {
	bi := big.NewInt(0).SetBytes(value[:])
	return bi.String()
}

func (gethCupDataFetcher *Fetcher) FetchPepData(methodArg interface{}, blockNumber int64) (*Peek, error) {
	abiJSON := MedianizerABI
	pepAddress := PepAddress
	method := "peek"
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	var r = &[]interface{}{
		ret0,
		ret1,
	}
	err := gethCupDataFetcher.blockchain.FetchContractData(abiJSON, pepAddress, method, methodArg, r, blockNumber)
	result := newResult(*ret0, *ret1)
	return result, err
}

func (gethCupDataFetcher *Fetcher) FetchPipData(methodArg interface{}, blockNumber int64) (*Peek, error) {
	abiJSON := MedianizerABI
	pipAddress := PipAddress
	method := "peek"
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	var r = &[]interface{}{
		ret0,
		ret1,
	}
	err := gethCupDataFetcher.blockchain.FetchContractData(abiJSON, pipAddress, method, methodArg, r, blockNumber)
	result := newResult(*ret0, *ret1)
	return result, err
}

type Per struct {
	Value *big.Int
}

func (per Per) Ray() string {
	return utils.Convert("ray", per.Value.String(), 17)
}

func (gethCupDataFetcher *Fetcher) FetchPerData(methodArg interface{}, blockNumber int64) (*Per, error) {
	abiJSON := tub.TubContractABI
	perAddress := PerAddress
	method := "per"
	var result = new(big.Int)
	err := gethCupDataFetcher.blockchain.FetchContractData(abiJSON, perAddress, method, methodArg, &result, blockNumber)
	return &Per{Value: result}, err
}
