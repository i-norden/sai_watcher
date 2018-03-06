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
	GetAllRows() ([]Row, error)
	MissingBlocks(startingBlockNumber int64, highestBlockNumber int64) ([]int64, error)
}

func (ebds DataStore) MissingBlocks(startingBlockNumber int64, highestBlockNumber int64) ([]int64, error) {
	// blocks that exist in vulcanize but no pep for
	numbers := make([]int64, 0)
	err := ebds.DB.Select(&numbers,
		`SELECT number
                    FROM blocks
                      LEFT JOIN maker.peps_everyblock
                        ON blocks.number = block_number
                    WHERE block_number ISNULL
                    AND number > $1
                    AND number <= $2
                LIMIT 20;`,
		startingBlockNumber,
		highestBlockNumber)
	if err != nil {
		return []int64{}, err
	}
	return numbers, err
}

type blockMetaData struct {
	BlockID   int   `db:"id"`
	BlockTime int64 `db:"time"`
}

func (ebds DataStore) Create(blockNumber int64, pep Peek, pip Peek, per Per) error {
	bmd := blockMetaData{}
	err := ebds.DB.Get(&bmd, `SELECT id, time FROM blocks WHERE number = $1`, blockNumber)
	if err != nil {
		return err
	}
	_, err = ebds.Exec(
		`INSERT INTO maker.peps_everyblock (block_number, pep, pip, per, block_id, block_time) 
                VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6)`,
		blockNumber, pep.Wad(), pip.Wad(), per.Ray(), bmd.BlockID, bmd.BlockTime)
	if err != nil {
		return err
	}
	return nil
}

func (ebds DataStore) Get(blockNumber int64) (*Row, error) {
	result := &Row{}
	err := ebds.DB.Get(result,
		`SELECT id, block_number, block_time, pep, pip, per 
                FROM maker.peps_everyblock WHERE block_number = $1`, blockNumber)
	if err != nil {
		return &Row{}, err
	}
	return result, nil
}

func (ebds DataStore) GetAllRows() ([]Row, error) {
	var results []Row
	err := ebds.DB.Select(&results, `SELECT id, block_number, block_id, block_time, pep, pip, per from maker.peps_everyblock`)
	if err != nil {
		return results, err
	}
	return results, nil
}
