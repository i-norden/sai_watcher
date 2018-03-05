package handlers

import (
	"log"

	"github.com/8thlight/sai_watcher/cup/fetchers"
	cuprepo "github.com/8thlight/sai_watcher/cup/repositories"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type CreatedHandler struct {
	core.ContractDataFetcher
	core.WatchedEvent
	shared.Handler
	cuprepo.ICupsRepository
	repositories.WatchedEventRepository
}

func NewCreatedHandler(db *postgres.DB, blockchain core.ContractDataFetcher) shared.Handler {
	var updater shared.Handler
	cr := cuprepo.CupsRepository{DB: db}
	we := repositories.WatchedEventRepository{DB: db}
	updater = &CreatedHandler{
		ICupsRepository:        cr,
		ContractDataFetcher:    blockchain,
		WatchedEventRepository: we,
	}
	fr := repositories.FilterRepository{DB: db}
	for _, filter := range createdFilters {
		fr.CreateFilter(filter)
	}
	return updater
}

var createdFilters = []filters.LogFilter{
	{
		Name:      "CupCreated",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0x89b8893b806db50897c8e2362c71571cfaeb9761ee40727f683f1793cda9df16"},
	},
}

func (u CreatedHandler) Execute() error {
	for _, filter := range createdFilters {
		watchedEvents, err := u.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println("Error fetching events for filter: ", filter.Name)
		}
		for _, watchedEvent := range watchedEvents {
			fetcher := fetchers.NewFetcher(u.ContractDataFetcher)
			args := common.HexToHash(watchedEvent.Data)
			cup, err := fetcher.FetchCupData(args, watchedEvent.BlockNumber)
			if err != nil {
				log.Println("Error fetching data for cup: ", watchedEvent.Data)
				log.Println(err)
			}
			cupsIndex := shared.HexToInt64(watchedEvent.Data)
			err = u.ICupsRepository.CreateCup(watchedEvent.LogID, *cup, watchedEvent.BlockNumber, false, cupsIndex)
			if err != nil {
				log.Println("Error persisting data for event: ", watchedEvent.LogID)
			}
		}
	}
	return nil
}
