package test_helpers

import "github.com/8thlight/sai_watcher/everyblock"

type MockEveryBlockDataStore struct {
	peps              []everyblock.Peek
	pips              []everyblock.Peek
	pers              []everyblock.Per
	BlockNumbers      []int64
	MissingBlocksData []int64
}

func (mpr *MockEveryBlockDataStore) GetAllRows() ([]everyblock.Row, error) {
	panic("implement me")
}

func (mpr *MockEveryBlockDataStore) MissingBlocks(startingBlockNumber int64, highestBlockNumber int64) ([]int64, error) {
	return mpr.MissingBlocksData, nil
}

func (mpr *MockEveryBlockDataStore) Get(blockNumber int64) (*everyblock.Row, error) {
	panic("implement me")
}

func (mpr *MockEveryBlockDataStore) Create(blockNumber int64, pep everyblock.Peek, pip everyblock.Peek, per everyblock.Per) error {
	mpr.BlockNumbers = append(mpr.BlockNumbers, blockNumber)
	mpr.peps = append(mpr.peps, pep)
	mpr.pips = append(mpr.pips, pip)
	mpr.pers = append(mpr.pers, per)
	return nil
}
