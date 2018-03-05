package repositories

import (
	"github.com/8thlight/sai_watcher/cup/fetchers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type ICupsRepository interface {
	CreateCup(logIndex int64, cup fetchers.Cup, blockNumber int64, isClosed bool, cupIndex int64) error
	GetCupEvents() ([]*core.WatchedEvent, error)
	GetCupsByIndex(cupIndex int) ([]DBCup, error)
}
