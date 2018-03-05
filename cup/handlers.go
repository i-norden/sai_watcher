package cup

import (
	"github.com/8thlight/sai_watcher/cup/handlers"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func HandlerInitializers() []shared.HandlerInitializer {
	return []shared.HandlerInitializer{
		handlers.NewCreatedHandler,
		handlers.NewShutHandler,
		handlers.NewDefaultHandler,
	}
}
