package everyblock_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/everyblock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Peps Repository", func() {
	var db *postgres.DB
	var pepsRepository everyblock.DataStore
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
		pepsRepository = everyblock.DataStore{DB: db}
		logsRepository = repositories.LogRepository{DB: db}
		filterRepository = repositories.FilterRepository{DB: db}
	})

	Describe("Creating a new pep record", func() {
		It("inserts new pep peek result with data", func() {
			wei := big.NewInt(0)
			wei.SetString("1000000000000000000", 10)
			ray := big.NewInt(0)
			ray.SetString("10000000000000000000000000000", 10)
			pip := everyblock.Peek{
				Value: everyblock.Value{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				OK:    false,
			}
			pep := everyblock.Peek{
				Value: everyblock.Value{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				OK:    false,
			}
			per := everyblock.Per{Value: ray}
			err = pepsRepository.Create(10, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())

			result, err := pepsRepository.Get(int64(10))

			Expect(err).ToNot(HaveOccurred())
			Expect(result.ID).ToNot(BeNil())
			Expect(result.Pep).To(Equal(big.NewRat(1, 1e18).FloatString(18)))
			Expect(result.Pip).To(Equal(big.NewRat(2, 1e18).FloatString(18)))
			Expect(result.Per).To(Equal("10"))

			Expect(result.BlockNumber).To(Equal(int64(10)))
		})
		It("returns values that do not have a record within a block range", func() {
			pip := everyblock.Peek{}
			pep := everyblock.Peek{}
			per := everyblock.Per{}
			err = pepsRepository.Create(0, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())
			err = pepsRepository.Create(1, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())
			err = pepsRepository.Create(2, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())
			err = pepsRepository.Create(3, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())
			err = pepsRepository.Create(4, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())

			result, err := pepsRepository.MissingBlocks(int64(0), int64(10))
			Expect(err).NotTo(HaveOccurred())

			Expect(err).ToNot(HaveOccurred())
			missingBlocks := []int64{5, 6, 7, 8, 9, 10}
			Expect(result).To(Equal(missingBlocks))
		})
	})
})
