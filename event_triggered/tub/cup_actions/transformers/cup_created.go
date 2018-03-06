package transformers

import (
	"log"

	"strings"

	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var CupCreatedFilter = filters.LogFilter{
	Name:      "CupCreated",
	FromBlock: 0,
	ToBlock:   -1,
	Address:   "0x448a5065aebb8e423f0896e6c5d525c040f59af3",
	// web3.utils.sha3('LogNewCup(address,bytes32)')
	Topics: core.Topics{"0x89b8893b806db50897c8e2362c71571cfaeb9761ee40727f683f1793cda9df16"},
}

type CupCreatedTransformer struct {
	Blockchain             core.Blockchain
	WatchedEventRepository datastore.WatchedEventRepository
	CupActionsRepository   cup_actions.CupActionsRepositoryInterface
}

func NewCupCreatedTransformer(db *postgres.DB, blockchain core.Blockchain) shared.Transformer {
	var transformer shared.Transformer
	we := repositories.WatchedEventRepository{DB: db}
	car := cup_actions.CupActionsRepository{DB: db}
	transformer = &CupCreatedTransformer{
		Blockchain:             blockchain,
		WatchedEventRepository: we,
		CupActionsRepository:   car,
	}
	fr := repositories.FilterRepository{DB: db}
	fr.CreateFilter(CupCreatedFilter)
	return transformer
}

func (cch CupCreatedTransformer) Execute() error {
	watchedEvents, err := cch.WatchedEventRepository.GetWatchedEvents(CupCreatedFilter.Name)
	if err != nil {
		log.Println("Error fetching events for filter: ", err)
	}
	for _, watchedEvent := range watchedEvents {
		model := models.CupAction{
			ID:              shared.HexToInt64(watchedEvent.Data),
			TransactionHash: watchedEvent.TxHash,
			Act:             "open",
			Arg:             "",
			Lad:             strings.ToLower(common.HexToAddress(watchedEvent.Topic1).Hex()),
			Ink:             "0",
			Art:             "0",
			Ire:             "0",
			Block:           watchedEvent.BlockNumber,
			Deleted:         false,
			Guy:             strings.ToLower(common.HexToAddress(watchedEvent.Topic1).Hex()),
		}
		cch.CupActionsRepository.CreateCupAction(model, watchedEvent.LogID)
	}
	return nil
}
