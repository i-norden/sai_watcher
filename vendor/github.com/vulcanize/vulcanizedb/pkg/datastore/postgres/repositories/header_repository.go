package repositories

import (
	"database/sql"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type HeaderRepository struct {
	database *postgres.DB
}

func NewHeaderRepository(database *postgres.DB) HeaderRepository {
	return HeaderRepository{database: database}
}

func (repository HeaderRepository) CreateOrUpdateHeader(header core.Header) (int64, error) {
	hash, err := repository.getHeaderHash(header)
	if err != nil {
		if headerDoesNotExist(err) {
			return repository.insertHeader(header)
		}
		return 0, err
	}
	if headerMustBeReplaced(hash, header) {
		return repository.replaceHeader(header)
	}
	return 0, err
}

func (repository HeaderRepository) GetHeader(blockNumber int64) (core.Header, error) {
	var header core.Header
	err := repository.database.Get(&header, `SELECT block_number, hash, raw FROM headers WHERE block_number = $1 AND eth_node_fingerprint = $2`,
		blockNumber, repository.database.Node.ID)
	return header, err
}

func (repository HeaderRepository) MissingBlockNumbers(startingBlockNumber, endingBlockNumber int64, nodeID string) []int64 {
	numbers := make([]int64, 0)
	repository.database.Select(&numbers, `SELECT all_block_numbers
	  FROM (
		  SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
	  WHERE all_block_numbers NOT IN (
		  SELECT block_number FROM headers WHERE eth_node_fingerprint = $3
	  ) `,
		startingBlockNumber, endingBlockNumber, nodeID)
	return numbers
}

func headerMustBeReplaced(hash string, header core.Header) bool {
	return hash != header.Hash
}

func headerDoesNotExist(err error) bool {
	return err == sql.ErrNoRows
}

func (repository HeaderRepository) getHeaderHash(header core.Header) (string, error) {
	var hash string
	err := repository.database.Get(&hash, `SELECT hash FROM headers WHERE block_number = $1 AND eth_node_fingerprint = $2`,
		header.BlockNumber, repository.database.Node.ID)
	return hash, err
}

func (repository HeaderRepository) insertHeader(header core.Header) (int64, error) {
	var headerId int64
	err := repository.database.QueryRowx(
		`INSERT INTO public.headers (block_number, hash, raw, eth_node_id, eth_node_fingerprint) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		header.BlockNumber, header.Hash, header.Raw, repository.database.NodeID, repository.database.Node.ID).Scan(&headerId)
	return headerId, err
}

func (repository HeaderRepository) replaceHeader(header core.Header) (int64, error) {
	_, err := repository.database.Exec(`DELETE FROM headers WHERE block_number = $1 AND eth_node_fingerprint = $2`,
		header.BlockNumber, repository.database.Node.ID)
	if err != nil {
		return 0, err
	}
	return repository.insertHeader(header)
}
