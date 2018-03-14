package everyblock

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type DataStore struct {
	*postgres.DB
}

type Repository interface {
	Create(blockNumber int64, pep Peek, pip Peek, per Per) error
	Get(blockNumber int64) (*Row, error)
}

func (ebds DataStore) Create(blockNumber int64, pep Peek, pip Peek, per Per) error {
	_, err := ebds.Exec(
		`INSERT INTO maker.peps_everyblock (block_number, pep, pip, per) 
                VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC)`,
		blockNumber, pep.Wad(), pip.Wad(), per.Ray())
	if err != nil {
		return err
	}
	return nil
}

func (ebds DataStore) Get(blockNumber int64) (*Row, error) {
	result := &Row{}
	err := ebds.DB.Get(result,
		`SELECT id, block_number, pep, pip, per 
                FROM maker.peps_everyblock WHERE block_number = $1`, blockNumber)
	if err != nil {
		return &Row{}, err
	}
	return result, nil
}
