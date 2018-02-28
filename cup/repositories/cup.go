package repositories

type DBCup struct {
	LogID       int   `db:"log_id"`
	CupIndex    int64 `db:"cup_index"`
	BlockNumber int64 `db:"block_number"`
	Lad         string
	Ink         string
	Art         string
	Irk         string
	IsClosed    bool `db:"is_closed"`
}
