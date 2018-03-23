package cup_actions

type CupActionModel struct {
	ID              int64
	TransactionHash string `db:"tx"`
	Act             string
	Arg             string
	Lad             string
	Ink             string
	Art             string
	Ire             string
	Block           int64
	Deleted         bool
	Guy             string
}
