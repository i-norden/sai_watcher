package test_helpers

import "github.com/8thlight/sai_watcher/event_triggered/cup_actions"

type MockFetcher struct {
	BlockNumbers []int64
	ReturnVal    cup_actions.Cup
}

func (mf *MockFetcher) FetchCupData(methodArg interface{}, blockNumber int64) (*cup_actions.Cup, error) {
	mf.BlockNumbers = append(mf.BlockNumbers, blockNumber)
	return &mf.ReturnVal, nil
}
