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
	"log"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type Transformer struct {
	Repository
	core.BlockChain
}

var PepsFilter = filters.LogFilter{
	Name:      "PepsLogger",
	FromBlock: 4742900,
	ToBlock:   -1,
	Address:   "0x99041f808d598b782d5a3e498681c2452a31da08",
	Topics:    core.Topics{"0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"},
}

func NewPepTransformer(db *postgres.DB, blockchain core.BlockChain) shared.Transformer {
	var transformer shared.Transformer
	pr := DataStore{DB: db}
	transformer = &Transformer{pr, blockchain}
	return transformer
}

func (pepTransformer *Transformer) Execute() error {
	fetcher := NewFetcher(pepTransformer.BlockChain)
	lastBlock := fetcher.blockchain.LastBlock().Int64()
	var blockUpperBound int64
	switch PepsFilter.ToBlock {
	case -1:
		blockUpperBound = lastBlock
	default:
		blockUpperBound = PepsFilter.ToBlock
	}
	blocks, err := pepTransformer.MissingBlocks(PepsFilter.FromBlock, blockUpperBound)
	if len(blocks) == 0 {
		return nil
	}
	if err != nil {
		log.Println("error fetching missing blocks: ", err)
	}
	for _, block := range blocks {
		pep, err := fetcher.FetchPepData(nil, block)
		if err != nil {
			log.Println("error fetching pep: ", err)
			return err
		}
		pip, err := fetcher.FetchPipData(nil, block)
		if err != nil {
			log.Println("error fetching pip: ", err)
			return err
		}
		per, err := fetcher.FetchPerData(nil, block)
		if err != nil {
			log.Println("error fetching per: ", err)
			return err
		}
		err = pepTransformer.Create(block, *pep, *pip, *per)
		if err != nil {
			log.Println("error creating pep: ", err)
			return err
		}

	}
	return nil
}
