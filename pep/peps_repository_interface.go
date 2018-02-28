package pep

import "github.com/vulcanize/vulcanizedb/pkg/core"

type IPepsRepository interface {
	CheckNewPep() ([]*core.WatchedEvent, error)
	CreatePep(value string, blockNumber int64, logId int64) error
}
