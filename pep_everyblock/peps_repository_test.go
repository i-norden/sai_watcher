package peps_everyblock_test

import (
	"github.com/8thlight/sai_watcher/pep_everyblock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Peps Repository", func() {
	var db *postgres.DB
	var pepsRepository peps_everyblock.PepsRepository
	var logsRepository repositories.LogRepository
	var filterRepository repositories.FilterRepository
	var err error

	BeforeEach(func() {
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())
		db.Query(`DELETE FROM maker.cups`)
		db.Query(`DELETE FROM maker.peps`)
		db.Query(`DELETE FROM maker.peps_everyblock`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM log_filters`)
		pepsRepository = peps_everyblock.PepsRepository{DB: db}
		logsRepository = repositories.LogRepository{DB: db}
		filterRepository = repositories.FilterRepository{DB: db}
	})

	Describe("Creating a new pep record", func() {
		It("inserts new pep peek result with data", func() {
			pep := peps_everyblock.DBPeekResult{
				Value: "10",
				OK:    true,
			}
			err = pepsRepository.CreatePep(pep.Value, 10)
			Expect(err).ToNot(HaveOccurred())

			var id int64
			var logId int64
			var blockNumber int64
			err = pepsRepository.DB.QueryRow(`SELECT id, block_number FROM maker.peps_everyblock`).
				Scan(&id, &blockNumber)
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeNil())
			Expect(logId).ToNot(BeNil())
			Expect(blockNumber).To(Equal(int64(10)))
		})
	})
})
