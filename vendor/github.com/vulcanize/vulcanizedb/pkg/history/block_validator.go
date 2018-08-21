package history

import (
	"io"
	"text/template"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
)

type BlockValidator struct {
	blockchain            core.BlockChain
	blockRepository       datastore.BlockRepository
	windowSize            int
	parsedLoggingTemplate template.Template
}

func NewBlockValidator(blockchain core.BlockChain, blockRepository datastore.BlockRepository, windowSize int) *BlockValidator {
	return &BlockValidator{
		blockchain,
		blockRepository,
		windowSize,
		ParsedWindowTemplate,
	}
}

func (bv BlockValidator) ValidateBlocks() ValidationWindow {
	window := MakeValidationWindow(bv.blockchain, bv.windowSize)
	blockNumbers := MakeRange(window.LowerBound, window.UpperBound)
	RetrieveAndUpdateBlocks(bv.blockchain, bv.blockRepository, blockNumbers)
	lastBlock := bv.blockchain.LastBlock().Int64()
	bv.blockRepository.SetBlocksStatus(lastBlock)
	return window
}

func (bv BlockValidator) Log(out io.Writer, window ValidationWindow) {
	bv.parsedLoggingTemplate.Execute(out, window)
}