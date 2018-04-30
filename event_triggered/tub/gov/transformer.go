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
	"fmt"

	"log"

	"github.com/8thlight/sai_watcher/event_triggered/tub"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var MoldActionHex = "0x92b0d72100000000000000000000000000000000000000000000000000000000"

var GovFilter = filters.LogFilter{
	Name:      "GovMold",
	FromBlock: 0,
	ToBlock:   -1,
	Address:   tub.TubContractAddress,
	Topics:    core.Topics{MoldActionHex},
}

func NewGovTransformer(db *postgres.DB, blockchain core.Blockchain) shared.Transformer {
	var transformer shared.Transformer
	gr := DataStore{DB: db}
	wer := repositories.WatchedEventRepository{DB: db}
	gf := GovFetcher{blockchain}
	transformer = &GovTransformer{
		Blockchain:             blockchain,
		WatchedEventRepository: wer,
		Fetcher:                gf,
		GovRepository:          gr,
	}
	fr := repositories.FilterRepository{DB: db}
	fr.CreateFilter(GovFilter)
	return transformer
}

type GovTransformer struct {
	Blockchain             core.Blockchain
	WatchedEventRepository datastore.WatchedEventRepository
	Fetcher                GovFetcherInterface
	GovRepository          Repository
}

func (gt GovTransformer) Execute() error {
	watchedEvents, err := gt.WatchedEventRepository.GetWatchedEvents(GovFilter.Name)
	if err != nil {
		fmt.Println("Error fetching watched events for gov: ", err)
	}
	for _, watchedEvent := range watchedEvents {
		blockNumber := watchedEvent.BlockNumber
		axe, err := gt.Fetcher.FetchAxe(blockNumber)
		if err != nil {
			log.Println(err)
		}
		cap, err := gt.Fetcher.FetchCap(blockNumber)
		if err != nil {
			log.Println(err)
		}
		fee, err := gt.Fetcher.FetchFee(blockNumber)
		if err != nil {
			log.Println(err)
		}
		gap, err := gt.Fetcher.FetchGap(blockNumber)
		if err != nil {
			log.Println(err)
		}
		mat, err := gt.Fetcher.FetchMat(blockNumber)
		if err != nil {
			log.Println(err)
		}
		tax, err := gt.Fetcher.FetchTax(blockNumber)
		ge := GovEntity{
			Block: watchedEvent.BlockNumber,
			Tx:    watchedEvent.TxHash,
			Var:   watchedEvent.Topic2,
			Arg:   watchedEvent.Topic3,
			Guy:   watchedEvent.Topic1,
			Cap:   cap.String(),
			Mat:   mat.String(),
			Tax:   tax.String(),
			Fee:   fee.String(),
			Axe:   axe.String(),
			Gap:   gap.String(),
		}
		gm := ConvertToModel(&ge)
		err = gt.GovRepository.CreateGov(gm, watchedEvent.LogID)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}
