package cup_actions

type CupActionEntity struct {
	ID              string
	TransactionHash string
	Act             string
	Arg             string
	Lad             string
	Ink             string
	Art             string
	Ire             string
	Block           int64
	Deleted         bool
}
