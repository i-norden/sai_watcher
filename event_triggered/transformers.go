package event_triggered

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/transformers"
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func TransformerInitializers() []shared.TransformerInitializer {
	return []shared.TransformerInitializer{
		transformers.NewCupCreatedTransformer,
		transformers.NewCupModifiedTransformer,
		gov.NewGovTransformer,
	}
}
