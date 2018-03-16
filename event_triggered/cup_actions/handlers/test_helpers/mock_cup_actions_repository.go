package test_helpers

import "github.com/8thlight/sai_watcher/event_triggered/cup_actions"

type MockCupActionsRepository struct {
	CupActions []cup_actions.CupActionModel
	LogIDs     []int64
}

func (car *MockCupActionsRepository) CreateCupAction(cupAction cup_actions.CupActionModel, logID int64) error {
	car.CupActions = append(car.CupActions, cupAction)
	car.LogIDs = append(car.LogIDs, logID)
	return nil
}
