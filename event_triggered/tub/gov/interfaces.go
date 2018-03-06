package gov

import "math/big"

type GovFetcherInterface interface {
	FetchCap(blockNumber int64) (*big.Int, error)
	FetchMat(blockNumber int64) (*big.Int, error)
	FetchTax(blockNumber int64) (*big.Int, error)
	FetchFee(blockNumber int64) (*big.Int, error)
	FetchAxe(blockNumber int64) (*big.Int, error)
	FetchGap(blockNumber int64) (*big.Int, error)
}

type Repository interface {
	CreateGov(govModel *GovModel, logID int64) error
	GetAllGovData() ([]GovModel, error)
}
