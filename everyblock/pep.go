package everyblock

type Row struct {
	ID          int64
	Pep         string
	Pip         string
	Per         string
	BlockNumber int64 `db:"block_number"`
	BlockID     int64 `db:"block_id"`
}
