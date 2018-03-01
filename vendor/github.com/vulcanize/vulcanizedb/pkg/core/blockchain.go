package core

import "math/big"

type Blockchain interface {
	CallContract(contractHash string, input []byte, blockNumber *big.Int) ([]byte, error)
	GetAttribute(contract Contract, attributeName string, blockNumber *big.Int) (interface{}, error)
	GetAttributes(contract Contract) (ContractAttributes, error)
	GetBlockByNumber(blockNumber int64) Block
	GetLogs(contract Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]Log, error)
	LastBlock() *big.Int
	Node() Node
}

type ContractDataFetcher interface {
	FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error
	GetContractOutput(address string, input []byte, blockNumber int64) ([]byte, error)
}
