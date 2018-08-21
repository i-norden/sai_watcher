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

package gov

import (
	"log"
	"math/big"

	"github.com/8thlight/sai_watcher/event_triggered/tub"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type GovFetcher struct {
	Blockchain core.BlockChain
}

func (gf GovFetcher) FetchCap(blockNumber int64) (*big.Int, error) {
	return gf.fetchGovData(blockNumber, "cap")
}

func (gf GovFetcher) FetchMat(blockNumber int64) (*big.Int, error) {
	return gf.fetchGovData(blockNumber, "mat")
}

func (gf GovFetcher) FetchTax(blockNumber int64) (*big.Int, error) {
	return gf.fetchGovData(blockNumber, "tax")
}

func (gf GovFetcher) FetchFee(blockNumber int64) (*big.Int, error) {
	return gf.fetchGovData(blockNumber, "fee")
}

func (gf GovFetcher) FetchAxe(blockNumber int64) (*big.Int, error) {
	return gf.fetchGovData(blockNumber, "axe")
}

func (gf GovFetcher) FetchGap(blockNumber int64) (*big.Int, error) {
	return gf.fetchGovData(blockNumber, "gap")
}

func (gf GovFetcher) fetchGovData(blockNumber int64, method string) (*big.Int, error) {
	var result = new(big.Int)
	err := gf.Blockchain.FetchContractData(tub.TubContractABI, tub.TubContractAddress, method, nil, &result, blockNumber)
	if err != nil {
		log.Printf("Error fetching %s: %s\n", method, err)
		return nil, err
	}
	return result, nil
}
