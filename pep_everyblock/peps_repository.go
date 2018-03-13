package peps_everyblock

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type PepsRepository struct {
	*postgres.DB
}

func (pepsRepository PepsRepository) CreatePep(value string, blockNumber int64) error {
	_, err := pepsRepository.Exec(
		`INSERT INTO maker.peps_everyblock (block_number, value) VALUES($1, $2)`, blockNumber, value)
	if err != nil {
		return err
	}
	return nil
}
