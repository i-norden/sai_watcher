package cup_actions

type CupActionsRepositoryInterface interface {
	CreateCupAction(cupAction CupActionModel, logID int64) error
}

type CupFetcherInterface interface {
	FetchCupData(methodArg interface{}, blockNumber int64) (*Cup, error)
}
