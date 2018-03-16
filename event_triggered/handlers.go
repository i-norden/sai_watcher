package event_triggered

import (
	"github.com/8thlight/sai_watcher/event_triggered/cup_actions/handlers"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func HandlerInitializers() []shared.HandlerInitializer {
	return []shared.HandlerInitializer{
		handlers.NewCupCreatedHandler,
		handlers.NewCupModifiedHandler,
	}
}
