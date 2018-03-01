package repositories

import (
	"github.com/8thlight/sai_watcher/cup/fetchers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type CupsRepository struct {
	*postgres.DB
}

func (cupsRepository CupsRepository) CreateCup(logId int64, cup fetchers.Cup, blockNumber int64, isClosed bool, cupIndex int64) error {
	_, err := cupsRepository.Exec(
		`INSERT INTO maker.cups (log_id, cup_index, lad, ink, art, irk, block_number, is_closed)
                SELECT $1, $2, $3, $4::NUMERIC, $5::NUMERIC, $6::NUMERIC, $7, $8
                WHERE NOT EXISTS (SELECT log_id FROM maker.cups WHERE log_id = $1)`,
		logId, cupIndex, cup.Lad.Hex(), cup.Ink.String(), cup.Art.String(), cup.Irk.String(), blockNumber, isClosed)
	if err != nil {
		return err
	}
	return nil
}

func (cupsRepository CupsRepository) GetCupsByIndex(cupIndex int) ([]DBCup, error) {
	cups := []DBCup{}

	err := cupsRepository.Select(&cups, "SELECT DISTINCT ON(cup_index, block_number) log_id,  cup_index, lad, ink::VARCHAR, art::VARCHAR, irk::VARCHAR, block_number, is_closed FROM maker.cups WHERE cup_index = $1", cupIndex)
	if err != nil {
		return cups, err
	}

	return cups, nil
}

func (cupsRepository CupsRepository) GetCupEvents() ([]*core.WatchedEvent, error) {
	rows, err := cupsRepository.DB.Queryx(`SELECT
              name, 
              watched_event_logs.id, 
              watched_event_logs.block_number, 
              address, 
              tx_hash, 
              index, 
              topic0, 
              topic1, 
              topic2, 
              topic3, 
              data 
            FROM watched_event_logs
              LEFT JOIN maker.cups cups
                ON cups.log_id = watched_event_logs.id
            WHERE cups.log_id IS NULL AND name = 'Cups'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lgs := make([]*core.WatchedEvent, 0)
	for rows.Next() {
		lg := new(core.WatchedEvent)
		err = rows.StructScan(lg)
		if err != nil {
			return nil, err
		}
		lgs = append(lgs, lg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lgs, nil
}
