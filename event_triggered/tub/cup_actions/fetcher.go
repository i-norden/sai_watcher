package cup_actions

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type CupFetcher struct {
	Blockchain core.Blockchain
}

var CupsContractMethod = "cups"

func (cupDataFetcher CupFetcher) FetchCupData(methodArg interface{}, blockNumber int64) (*Cup, error) {
	abiJSON := tub.TubContractABI
	address := tub.TubContractAddress
	method := CupsContractMethod
	result := &Cup{}
	err := cupDataFetcher.Blockchain.FetchContractData(abiJSON, address, method, methodArg, result, blockNumber)
	return result, err
}
