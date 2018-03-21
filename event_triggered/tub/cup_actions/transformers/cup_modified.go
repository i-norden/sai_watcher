package transformers

import (
	"log"

	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type CupModifiedTransformer struct {
	CupActionsRepository   cup_actions.CupActionsRepositoryInterface
	WatchedEventRepository datastore.WatchedEventRepository
	FilterRepository       datastore.FilterRepository
	Fetcher                cup_actions.CupFetcherInterface
	Blockchain             core.Blockchain
}

var GiveActionHex = "0xbaa8529c00000000000000000000000000000000000000000000000000000000"
var LockActionHex = "0xb3b77a5100000000000000000000000000000000000000000000000000000000"
var FreeActionHex = "0xa5cd184e00000000000000000000000000000000000000000000000000000000"
var DrawActionHex = "0x440f19ba00000000000000000000000000000000000000000000000000000000"
var WipeActionHex = "0x73b3810100000000000000000000000000000000000000000000000000000000"
var ShutActionHex = "0xb84d210600000000000000000000000000000000000000000000000000000000"
var BiteActionHex = "0x40cc885400000000000000000000000000000000000000000000000000000000"
var TubContractAddress = "0x448a5065aebb8e423f0896e6c5d525c040f59af3"

var CupModifiedFilters = []filters.LogFilter{
	{
		Name:      "CupsGive",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{GiveActionHex}},
	{
		Name:      "CupsLock",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{LockActionHex},
	},
	{
		Name:      "CupsFree",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{FreeActionHex},
	},
	{
		Name:      "CupsDraw",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{DrawActionHex},
	},
	{
		Name:      "CupsWipe",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{WipeActionHex},
	},
	{
		Name:      "CupsBite",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{BiteActionHex},
	},
	{
		Name:      "CupShut",
		FromBlock: 0,
		ToBlock:   -1,
		Address:   TubContractAddress,
		Topics:    core.Topics{ShutActionHex},
	},
}

func NewCupModifiedTransformer(db *postgres.DB, blockchain core.Blockchain) shared.Transformer {
	car := cup_actions.CupActionsRepository{DB: db}
	fr := repositories.FilterRepository{DB: db}
	wer := repositories.WatchedEventRepository{DB: db}
	fetcher := cup_actions.CupFetcher{blockchain}
	for _, filter := range CupModifiedFilters {
		fr.CreateFilter(filter)
	}
	return &CupModifiedTransformer{
		CupActionsRepository:   car,
		Blockchain:             blockchain,
		FilterRepository:       fr,
		WatchedEventRepository: wer,
		Fetcher:                fetcher,
	}
}

func (cmh CupModifiedTransformer) Execute() error {
	for _, filter := range CupModifiedFilters {
		watchedEvents, err := cmh.WatchedEventRepository.GetWatchedEvents(filter.Name)
		if err != nil {
			log.Println("Error fetching events for filter: ", filter.Name)
		}
		for _, watchedEvent := range watchedEvents {
			createCupAction(cmh, watchedEvent)
		}
	}
	return nil
}

func createCupAction(cmh CupModifiedTransformer, watchedEvent *core.WatchedEvent) {
	args := common.HexToHash(watchedEvent.Topic2)
	cup, err := cmh.Fetcher.FetchCupData(args, watchedEvent.BlockNumber)
	if err != nil {
		log.Println("Error fetching cup data: ", err)
	}
	entity := cup_actions.CupActionEntity{
		ID:              watchedEvent.Topic2,
		TransactionHash: watchedEvent.TxHash,
		Act:             getCupActionName(watchedEvent.Topic0),
		Arg:             watchedEvent.Topic3,
		Lad:             watchedEvent.Topic1,
		Ink:             cup.Ink.String(),
		Art:             cup.Art.String(),
		Ire:             cup.Ire.String(),
		Block:           watchedEvent.BlockNumber,
		Deleted:         isShut(watchedEvent.Topic0),
	}
	model := cup_actions.ConvertToModel(entity)
	cmh.CupActionsRepository.CreateCupAction(model, watchedEvent.LogID)
}

func getCupActionName(act string) string {
	actMappings := map[string]string{
		GiveActionHex: "give",
		LockActionHex: "lock",
		FreeActionHex: "free",
		DrawActionHex: "draw",
		WipeActionHex: "wipe",
		ShutActionHex: "shut",
		BiteActionHex: "bite",
	}
	return actMappings[act]
}

func isShut(act string) bool {
	return act == ShutActionHex
}
