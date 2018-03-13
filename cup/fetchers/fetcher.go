package fetchers

import (
	"errors"

	"math/big"

	"path/filepath"

	"github.com/8thlight/sai_watcher/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

var (
	ErrApiRequestFailed = errors.New("etherscan api request failed")
)

func NewFetcher(blockchain core.Blockchain) *Fetcher {
	return &Fetcher{
		blockchain: blockchain,
	}
}

type Fetcher struct {
	blockchain core.Blockchain
}

type Cup struct {
	Lad common.Address
	Ink *big.Int
	Art *big.Int
	Irk *big.Int
}

func (gethCupDateFetcher *Fetcher) FetchCupData(methodArg interface{}, blockNumber int64) (*Cup, error) {
	abiPath := filepath.Join(utils.ProjectRoot(), "cup", "saitub.json")
	abiJSON, err := geth.ReadAbiFile(abiPath)
	if err != nil {
		return &Cup{}, err
	}
	address := "0x448a5065aebb8e423f0896e6c5d525c040f59af3"
	method := "cups"
	result := &Cup{}
	err = gethCupDateFetcher.blockchain.FetchContractData(abiJSON, address, method, methodArg, result, blockNumber)
	return result, err
}
