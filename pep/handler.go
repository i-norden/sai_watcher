package pep

import (
	"log"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type Handler struct {
	IPepsRepository
	datastore.WatchedEventRepository
}

var pepsFilters = []filters.LogFilter{
	{
		Name:      "PepsLogger",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   "0x99041f808d598b782d5a3e498681c2452a31da08",
		Topics:    core.Topics{"0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"},
	},
}

func NewPepHandler(db *postgres.DB, blockchain core.ContractDataFetcher) shared.Handler {
	var handler shared.Handler
	pr := PepsRepository{DB: db}
	we := repositories.WatchedEventRepository{DB: db}
	handler = &Handler{pr, we}
	fr := repositories.FilterRepository{DB: db}
	for _, filter := range pepsFilters {
		fr.CreateFilter(filter)
	}
	return handler
}

func (pepHandler *Handler) Execute() error {
	for _, filter := range pepsFilters {
		watchedEvents, err := pepHandler.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println("Error fetching events for filter: ", filter.Name)
		}
		for _, watchedEvent := range watchedEvents {
			createPepData(watchedEvent, pepHandler)
		}
	}
	return nil
}

func createPepData(watchedEvent *core.WatchedEvent, pepHandler *Handler) {
	value := shared.HexToString(watchedEvent.Data)
	err := pepHandler.CreatePep(value, watchedEvent.BlockNumber, watchedEvent.LogID)
	if err != nil {
		log.Println("Error persisting data for event: ", watchedEvent.LogID)
	}
}
