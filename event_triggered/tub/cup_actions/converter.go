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

package cup_actions

import (
	"strings"

	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/models"
	"github.com/8thlight/sai_watcher/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func ConvertToModel(entity CupActionEntity) models.CupAction {
	return models.CupAction{
		ID:              shared.HexToInt64(entity.ID),
		TransactionHash: entity.TransactionHash,
		Act:             entity.Act,
		Arg:             Arg(entity.Arg, entity.Act),
		Lad:             strings.ToLower(common.HexToAddress(entity.Lad).Hex()),
		Ink:             utils.Convert("wad", entity.Ink, 17),
		Art:             utils.Convert("wad", entity.Art, 17),
		Ire:             utils.Convert("wad", entity.Ire, 17),
		Block:           entity.Block,
		Deleted:         entity.Deleted,
		Guy:             strings.ToLower(entity.Guy),
	}
}

func Arg(s string, act string) string {
	if act == "give" {
		return s
	}
	return utils.Arg(s)
}
