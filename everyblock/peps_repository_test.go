package everyblock_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/everyblock"
	"github.com/8thlight/sai_watcher/test_helpers"
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
	var blockRepository repositories.BlockRepository
	var err error

	BeforeEach(func() {
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())
		db.Query(`DELETE FROM maker.peps_everyblock`)
		db.Query(`DELETE FROM blocks`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM log_filters`)
		pepsRepository = everyblock.DataStore{DB: db}
		blockRepository = repositories.BlockRepository{DB: db}
	})

	Describe("Creating a new pep record", func() {
		It("inserts new pep peek result with data", func() {
			err := blockRepository.CreateOrUpdateBlock(core.Block{Number: 10, Time: int64(100)})
			Expect(err).ToNot(HaveOccurred())
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
			Expect(result.BlockTime).To(Equal(int64(100)))
		})
	})

	Describe("Handling reorgs", func() {
		It("removes pep peek result on block reorg", func() {
			blockNumber := int64(12345678)
			block := core.Block{
				Number:       blockNumber,
				Transactions: []core.Transaction{{}},
			}
			err := blockRepository.CreateOrUpdateBlock(block)
			Expect(err).ToNot(HaveOccurred())
			var blockID int64
			err = blockRepository.Get(&blockID, `Select id from blocks`)
			Expect(err).NotTo(HaveOccurred())

			// confirm newly created Pep is present with existing block ID
			err = pepsRepository.Create(blockNumber, everyblock.Peek{}, everyblock.Peek{}, everyblock.Per{})
			Expect(err).NotTo(HaveOccurred())
			result := &everyblock.Row{}
			err = pepsRepository.DB.QueryRowx(
				`SELECT * FROM maker.peps_everyblock WHERE block_id = $1`, blockID).StructScan(result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.BlockID).To(Equal(blockID))

			// block is removed because of reorg
			_, err = blockRepository.DB.Exec(`DELETE FROM blocks WHERE id = $1`, blockID)
			Expect(err).ToNot(HaveOccurred())
			var blockCount int
			err = blockRepository.Get(&blockCount, `SELECT count(*) FROM logs WHERE id = $1`, blockID)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockCount).To(BeZero())

			// confirm corresponding pep is removed
			var pepCount int
			err = pepsRepository.DB.QueryRowx(
				`SELECT count(*) FROM maker.peps_everyblock WHERE block_id = $1`, blockID).Scan(&pepCount)
			Expect(err).ToNot(HaveOccurred())
			Expect(pepCount).To(BeZero())
		})
	})

	Describe("Fetching missing blocks", func() {
		It("returns values that do not have a record AND are in vulcanize db within a block range", func() {
			for i := 0; i < 11; i++ {
				err = blockRepository.CreateOrUpdateBlock(core.Block{Number: int64(i)})
				Expect(err).ToNot(HaveOccurred())
			}

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

	Describe("Getting all rows", func() {
		It("returns data for every row", func() {
			blockNumberOne := int64(1)
			err := blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumberOne, Time: int64(100)})
			Expect(err).ToNot(HaveOccurred())
			rayOne := big.NewInt(0)
			rayOne.SetString("20000000000000000000000000000", 10)
			pipOne := everyblock.Peek{
				Value: everyblock.Value{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
				OK:    false,
			}
			pepOne := everyblock.Peek{
				Value: everyblock.Value{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
				OK:    false,
			}
			perOne := everyblock.Per{Value: rayOne}
			err = pepsRepository.Create(blockNumberOne, pepOne, pipOne, perOne)
			blockNumberTwo := int64(2)
			err = blockRepository.CreateOrUpdateBlock(core.Block{Number: blockNumberTwo, Time: int64(100)})
			Expect(err).ToNot(HaveOccurred())
			rayTwo := big.NewInt(0)
			rayTwo.SetString("30000000000000000000000000000", 10)
			pipTwo := everyblock.Peek{
				Value: everyblock.Value{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5},
				OK:    false,
			}
			pepTwo := everyblock.Peek{
				Value: everyblock.Value{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
				OK:    false,
			}
			perTwo := everyblock.Per{Value: rayTwo}
			err = pepsRepository.Create(blockNumberTwo, pepTwo, pipTwo, perTwo)

			results, err := pepsRepository.GetAllRows()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(2))

			rowOne := results[0]
			Expect(rowOne.BlockNumber).To(Equal(blockNumberOne))
			actualRowOnePip := test_helpers.GetFloat(rowOne.Pip)
			Expect(actualRowOnePip.String()).To(Equal(pipOne.Wad()))
			actualRowOnePep := test_helpers.GetFloat(rowOne.Pep)
			Expect(actualRowOnePep.String()).To(Equal(pepOne.Wad()))
			actualRowOnePer := test_helpers.GetFloat(rowOne.Per)
			Expect(actualRowOnePer.String()).To(Equal(perOne.Ray()))

			rowTwo := results[1]
			Expect(rowTwo.BlockNumber).To(Equal(blockNumberTwo))
			actualRowTwoPip := test_helpers.GetFloat(rowTwo.Pip)
			Expect(actualRowTwoPip.String()).To(Equal(pipTwo.Wad()))
			actualRowTwoPep := test_helpers.GetFloat(rowTwo.Pep)
			Expect(actualRowTwoPep.String()).To(Equal(pepTwo.Wad()))
			actualRowTwoPer := test_helpers.GetFloat(rowTwo.Per)
			Expect(actualRowTwoPer.String()).To(Equal(perTwo.Ray()))
		})
	})
})
