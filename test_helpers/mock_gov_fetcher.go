package test_helpers

import "math/big"

type MockGovFetcher struct {
	CapCalls []int64
	MatCalls []int64
	TaxCalls []int64
	FeeCalls []int64
	AxeCalls []int64
	GapCalls []int64
}

func (mgf *MockGovFetcher) FetchAxe(blockNumber int64) (*big.Int, error) {
	mgf.AxeCalls = append(mgf.AxeCalls, blockNumber)
	return big.NewInt(123), nil
}

func (mgf *MockGovFetcher) FetchCap(blockNumber int64) (*big.Int, error) {
	mgf.CapCalls = append(mgf.CapCalls, blockNumber)
	return big.NewInt(123), nil
}

func (mgf *MockGovFetcher) FetchFee(blockNumber int64) (*big.Int, error) {
	mgf.FeeCalls = append(mgf.FeeCalls, blockNumber)
	return big.NewInt(123), nil
}

func (mgf *MockGovFetcher) FetchGap(blockNumber int64) (*big.Int, error) {
	mgf.GapCalls = append(mgf.GapCalls, blockNumber)
	return big.NewInt(123), nil
}

func (mgf *MockGovFetcher) FetchMat(blockNumber int64) (*big.Int, error) {
	mgf.MatCalls = append(mgf.MatCalls, blockNumber)
	return big.NewInt(123), nil
}

func (mgf *MockGovFetcher) FetchTax(blockNumber int64) (*big.Int, error) {
	mgf.TaxCalls = append(mgf.TaxCalls, blockNumber)
	return big.NewInt(123), nil
}
