package everyblock

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared"
)

func TransformerInitializers() []shared.TransformerInitializer {
	return []shared.TransformerInitializer{
		NewPepTransformer,
	}
}
