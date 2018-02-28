package handlers

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	cuprepo "github.com/8thlight/sai_watcher/cup/repositories"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
	"github.com/8thlight/sai_watcher/cup/fetchers"
)

type ShutHandler struct {
	core.WatchedEvent
	shared.Handler
	core.ContractDataFetcher
	cuprepo.ICupsRepository
	repositories.WatchedEventRepository
}

var shutFilters = []filters.LogFilter{
	{
		Name:      "CupShut",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0xb84d210600000000000000000000000000000000000000000000000000000000"},
	},
}

func NewShutHandler(db *postgres.DB, blockchain core.ContractDataFetcher) shared.Handler {
	var updater shared.Handler
	cr := cuprepo.CupsRepository{DB: db}
	we := repositories.WatchedEventRepository{DB: db}
	updater = &ShutHandler{
		ICupsRepository:        cr,
		ContractDataFetcher:    blockchain,
		WatchedEventRepository: we,
	}
	fr := repositories.FilterRepository{DB: db}
	for _, filter := range shutFilters {
		fr.CreateFilter(filter)
	}
	return updater
}

func (u ShutHandler) Execute() error {
	for _, filter := range shutFilters {
		watchedEvents, err := u.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println("Error fetching events for filter: ", filter.Name)
		}
		for _, watchedEvent := range watchedEvents {
			shutCup(u, watchedEvent)
		}
	}
	return nil
}

func shutCup(u ShutHandler, watchedEvent *core.WatchedEvent) {
	fetcher := fetchers.NewFetcher(u.ContractDataFetcher)
	args := common.HexToHash(watchedEvent.Topic2)
	cup, err := fetcher.FetchCupData(args, watchedEvent.BlockNumber)
	if err != nil {
		log.Println("Error fetching data for cup: ", watchedEvent.Topic2)
		// TODO: Continue to next event on error?
	}
	cupsIndex := shared.HexToInt64(watchedEvent.Topic2)
	err = u.ICupsRepository.CreateCup(watchedEvent.LogID, *cup, watchedEvent.BlockNumber, true, cupsIndex)
	if err != nil {
		log.Println("Error handling event: ", watchedEvent)
	}
}
