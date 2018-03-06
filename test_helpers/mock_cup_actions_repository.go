package test_helpers

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/models"
)

type MockCupActionsRepository struct {
	CupActions []models.CupAction
	LogIDs     []int64
}

func (car *MockCupActionsRepository) GetAllCupData() ([]models.Cup, error) {
	panic("implement me")
}

func (car *MockCupActionsRepository) CreateCupAction(cupAction models.CupAction, logID int64) error {
	car.CupActions = append(car.CupActions, cupAction)
	car.LogIDs = append(car.LogIDs, logID)
	return nil
}
