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

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
)

type MockGovRepository struct {
	Govs   []gov.GovModel
	LogIDs []int64
}

func (ds MockGovRepository) GetAllGovData() ([]gov.GovModel, error) {
	panic("implement me")
}

func (ds MockGovRepository) CreateGov(govModel *gov.GovModel, logID int64) error {
	ds.Govs = append(ds.Govs, *govModel)
	ds.LogIDs = append(ds.LogIDs, logID)
	return nil
}
