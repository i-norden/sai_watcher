package cup_actions

import (
	"github.com/8thlight/sai_watcher/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func ConvertToModel(entity CupActionEntity) CupActionModel {
	return CupActionModel{
		ID:              shared.HexToInt64(entity.ID),
		TransactionHash: entity.TransactionHash,
		Act:             entity.Act,
		Arg:             utils.Arg(entity.Arg),
		Lad:             common.HexToAddress(entity.Lad).Hex(),
		Ink:             utils.Convert("wad", entity.Ink, 17),
		Art:             utils.Convert("wad", entity.Art, 17),
		Ire:             utils.Convert("wad", entity.Ire, 17),
		Block:           entity.Block,
		Deleted:         entity.Deleted,
	}
}
