package cup_actions

import "github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/models"

type CupActionsRepositoryInterface interface {
	CreateCupAction(cupAction models.CupAction, logID int64) error
	GetAllCupData() ([]models.Cup, error)
}

type CupFetcherInterface interface {
	FetchCupData(methodArg interface{}, blockNumber int64) (*Cup, error)
}
