package peps_everyblock

import (
	"path/filepath"

	"math/big"

	"github.com/8thlight/sai_watcher/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

func NewFetcher(blockchain core.Blockchain) *Fetcher {
	return &Fetcher{
		blockchain: blockchain,
	}
}

type Fetcher struct {
	blockchain core.Blockchain
}

type Value [32]byte

type PeekResult struct {
	Value
	OK bool
}

func newResult(value [32]byte, ok bool) *PeekResult {
	return &PeekResult{Value: value, OK: ok}
}

func (value Value) Hex() string {
	return hexutil.Encode(value[:])
}

func (value Value) String() string {
	bi := big.NewInt(0).SetBytes(value[:])
	return bi.String()
}

func (gethCupDateFetcher *Fetcher) FetchContractData(methodArg interface{}, blockNumber int64) (*PeekResult, error) {
	abiPath := filepath.Join(utils.ProjectRoot(), "pep_everyblock", "medianizer.json")
	abiJSON, err := geth.ReadAbiFile(abiPath)
	address := "0x99041f808d598b782d5a3e498681c2452a31da08"
	method := "peek"
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	var r = &[]interface{}{
		ret0,
		ret1,
	}
	err = gethCupDateFetcher.blockchain.FetchContractData(abiJSON, address, method, methodArg, r, blockNumber)
	result := newResult(*ret0, *ret1)
	return result, err
}
