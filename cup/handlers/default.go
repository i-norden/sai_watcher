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

type DefaultHandler struct {
	core.WatchedEvent
	core.Blockchain
	shared.Handler
	cuprepo.ICupsRepository
	repositories.FilterRepository
	repositories.WatchedEventRepository
}

var cupsFilters = []filters.LogFilter{
	{
		Name:      "CupsGive",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0xbaa8529c00000000000000000000000000000000000000000000000000000000"}},
	{
		Name:      "CupsLock",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0xb3b77a5100000000000000000000000000000000000000000000000000000000"},
	},
	{
		Name:      "CupsFree",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0xa5cd184e00000000000000000000000000000000000000000000000000000000"},
	},
	{
		Name:      "CupsDraw",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0x440f19ba00000000000000000000000000000000000000000000000000000000"},
	},
	{
		Name:      "CupsWipe",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0x73b3810100000000000000000000000000000000000000000000000000000000"},
	},
	{
		Name:      "CupsBite",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
		Topics:    core.Topics{"0x40cc885400000000000000000000000000000000000000000000000000000000"},
	},
}

func NewDefaultHandler(db *postgres.DB, blockchain core.Blockchain) shared.Handler {
	cr := cuprepo.CupsRepository{DB: db}
	fr := repositories.FilterRepository{DB: db}
	we := repositories.WatchedEventRepository{DB: db}
	for _, filter := range cupsFilters {
		fr.CreateFilter(filter)
	}
	return &DefaultHandler{
		ICupsRepository:        cr,
		Blockchain:             blockchain,
		FilterRepository:       fr,
		WatchedEventRepository: we,
	}
}

func (u DefaultHandler) Execute() error {
	for _, filter := range cupsFilters {
		watchedEvents, err := u.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println("Error fetching events for filter: ", filter.Name)
		}
		for _, watchedEvent := range watchedEvents {
			createCupData(u, watchedEvent)
		}
	}
	return nil
}

func createCupData(u DefaultHandler, watchedEvent *core.WatchedEvent) {
	fetcher := fetchers.NewFetcher(u.Blockchain)
	args := common.HexToHash(watchedEvent.Topic2)
	cup, err := fetcher.FetchCupData(args, watchedEvent.BlockNumber)
	if err != nil {
		log.Println("Error fetching data for cup: ", watchedEvent.Topic2)
	}
	cupsIndex := shared.HexToInt64(watchedEvent.Topic2)
	err = u.ICupsRepository.CreateCup(watchedEvent.LogID, *cup, watchedEvent.BlockNumber, false, cupsIndex)
	if err != nil {
		log.Println("Error persisting data for event: ", watchedEvent.LogID)
	}
}
