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
