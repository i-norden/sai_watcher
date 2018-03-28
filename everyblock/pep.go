package everyblock

type Row struct {
	Pep         string
	Pip         string
	Per         string
	BlockNumber int64 `db:"block_number"`
	BlockTime   int64 `db:"block_time"`
}
