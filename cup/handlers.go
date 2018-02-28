package cup

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/8thlight/sai_watcher/cup/handlers"
)

func HandlerInitializers() []shared.HandlerInitializer {
	return []shared.HandlerInitializer{
		handlers.NewCreatedHandler,
		handlers.NewShutHandler,
		handlers.NewDefaultHandler,
	}
}
