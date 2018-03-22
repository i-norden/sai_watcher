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
	MissingBlocks(startingBlockNumber int64, highestBlockNumber int64) ([]int64, error)
}

func (ebds DataStore) MissingBlocks(startingBlockNumber int64, highestBlockNumber int64) ([]int64, error) {
	numbers := make([]int64, 0)
	err := ebds.DB.Select(&numbers,
		`SELECT all_block_numbers
            FROM (
                SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
                LEFT JOIN maker.peps_everyblock
                    ON block_number = all_block_numbers
            WHERE block_number ISNULL
            Limit 20`,
		startingBlockNumber,
		highestBlockNumber)
	if err != nil {
		return []int64{}, err
	}
	return numbers, err
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
