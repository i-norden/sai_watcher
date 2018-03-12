package peps_everyblock

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func HandlerInitializers() []shared.HandlerInitializer {
	return []shared.HandlerInitializer{
		NewPepHandler,
	}
}
