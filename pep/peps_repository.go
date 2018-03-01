package pep

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type PepsRepository struct {
	*postgres.DB
}

func (pepsRepository PepsRepository) CreatePep(value string, blockNumber int64, logId int64) error {
	_, err := pepsRepository.Exec(
		`INSERT INTO maker.peps (log_id, block_number, value)
                SELECT $1, $2, $3 
                WHERE NOT EXISTS ( SELECT log_id FROM maker.peps WHERE log_id = $1)`, logId, blockNumber, value)
	if err != nil {
		return err
	}
	return nil
}

func (pepsRepository PepsRepository) CheckNewPep() ([]*core.WatchedEvent, error) {
	rows, err := pepsRepository.DB.Queryx(`SELECT
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
              LEFT JOIN maker.peps peps
                ON peps.log_id = watched_event_logs.id
            WHERE peps.log_id IS NULL AND name = 'Peps'`)
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
