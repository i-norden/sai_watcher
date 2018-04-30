// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
