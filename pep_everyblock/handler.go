package peps_everyblock

import (
	"log"

	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type Handler struct {
	IPepsRepository
	core.Blockchain
}

var pepsFilters = []filters.LogFilter{
	{
		Name:      "PepsLogger",
		FromBlock: 5209657,
		ToBlock:   -1,
		Address:   "0x99041f808d598b782d5a3e498681c2452a31da08",
		Topics:    core.Topics{"0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"},
	},
}

func NewPepHandler(db *postgres.DB, blockchain core.Blockchain) shared.Handler {
	var handler shared.Handler
	pr := PepsRepository{DB: db}
	handler = &Handler{pr, blockchain}
	fr := repositories.FilterRepository{DB: db}
	for _, filter := range pepsFilters {
		fr.CreateFilter(filter)
	}
	return handler
}

func (pepHandler *Handler) Execute() error {
	fetcher := NewFetcher(pepHandler.Blockchain)
	lastBlock := fetcher.blockchain.LastBlock().Int64()
	for _, filter := range pepsFilters {
		var blockUpperBound int64
		if filter.ToBlock == -1 {
			blockUpperBound = lastBlock
		}
		for i := filter.FromBlock; i <= blockUpperBound; i++ {
			pepValue, err := fetcher.FetchContractData(nil, i)
			if err != nil {
				log.Println("Error fetching events for filter: ", filter.Name)
				log.Println("Error: ", err)
			}
			err = pepHandler.CreatePep(pepValue.String(), i)
			if err != nil {
				log.Println("Error creating pep: ", err)
			}
		}

	}
	return nil
}
