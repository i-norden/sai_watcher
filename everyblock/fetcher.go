package everyblock

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

type Peek struct {
	Value
	OK bool
}

func newResult(value [32]byte, ok bool) *Peek {
	return &Peek{Value: value, OK: ok}
}

func (peek Peek) Wad() string {
	return utils.Convert("wad", peek.Value.String(), 15)
}

func (value Value) Hex() string {
	return hexutil.Encode(value[:])
}

func (value Value) String() string {
	bi := big.NewInt(0).SetBytes(value[:])
	return bi.String()
}

func (gethCupDateFetcher *Fetcher) FetchPepData(methodArg interface{}, blockNumber int64) (*Peek, error) {
	abiPath := filepath.Join(utils.ProjectRoot(), "everyblock", "medianizer.json")
	abiJSON, err := geth.ReadAbiFile(abiPath)
	pepAddress := "0x99041f808d598b782d5a3e498681c2452a31da08"
	method := "peek"
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	var r = &[]interface{}{
		ret0,
		ret1,
	}
	err = gethCupDateFetcher.blockchain.FetchContractData(abiJSON, pepAddress, method, methodArg, r, blockNumber)
	result := newResult(*ret0, *ret1)
	return result, err
}

func (gethCupDateFetcher *Fetcher) FetchPipData(methodArg interface{}, blockNumber int64) (*Peek, error) {
	abiPath := filepath.Join(utils.ProjectRoot(), "everyblock", "medianizer.json")
	abiJSON, err := geth.ReadAbiFile(abiPath)
	pipAddress := "0x729d19f657bd0614b4985cf1d82531c67569197b"
	method := "peek"
	var (
		ret0 = new([32]byte)
		ret1 = new(bool)
	)
	var r = &[]interface{}{
		ret0,
		ret1,
	}
	err = gethCupDateFetcher.blockchain.FetchContractData(abiJSON, pipAddress, method, methodArg, r, blockNumber)
	result := newResult(*ret0, *ret1)
	return result, err
}

type Per struct {
	Value *big.Int
}

func (per Per) Ray() string {
	return utils.Convert("ray", per.Value.String(), 17)
}

func (gethCupDateFetcher *Fetcher) FetchPerData(methodArg interface{}, blockNumber int64) (*Per, error) {
	abiPath := filepath.Join(utils.ProjectRoot(), "everyblock", "tub.json")
	abiJSON, err := geth.ReadAbiFile(abiPath)
	perAddress := "0x448a5065aebb8e423f0896e6c5d525c040f59af3"
	method := "per"
	var result = new(big.Int)
	err = gethCupDateFetcher.blockchain.FetchContractData(abiJSON, perAddress, method, methodArg, &result, blockNumber)
	return &Per{Value: result}, err
}
