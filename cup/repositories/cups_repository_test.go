package repositories

import (
	"math/big"

	"github.com/8thlight/sai_watcher/cup/fetchers"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var _ = Describe("Cups Repository", func() {
	var db *postgres.DB
	var cupsRepository CupsRepository
	var filterRepository repositories.FilterRepository
	var logRepository repositories.LogRepository
	var node core.Node

	BeforeEach(func() {
		node = core.Node{}
		var err error
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())
		db.Query(`DELETE FROM maker.cups`)
		db.Query(`DELETE FROM maker.peps`)
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM log_filters`)
		logRepository = repositories.LogRepository{DB: db}
		cupsRepository = CupsRepository{DB: db}
		filterRepository = repositories.FilterRepository{DB: db}
	})

	Describe("Creating a new cups record", func() {
		It("inserts a new cup", func() {
			err := logRepository.CreateLogs([]core.Log{{}})
			Expect(err).ToNot(HaveOccurred())

			cupIndex := int64(8675)
			lad := common.StringToAddress("Address")
			ink := big.NewInt(1)
			art := big.NewInt(29259384232352)
			irk := big.NewInt(3)
			cup := fetchers.Cup{
				Lad: lad,
				Ink: ink,
				Art: art,
				Irk: irk,
			}

			var logId int64
			err = logRepository.Get(&logId, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())

			blockNumber := int64(1234)
			err = cupsRepository.CreateCup(logId, cup, blockNumber, false, cupIndex)
			Expect(err).ToNot(HaveOccurred())

			var DBlogId int64
			var DBcupIndex int64
			var DBlad string
			var DBink int64
			var DBart int64
			var DBirk int64
			var DBblockNumber int64
			var DBisClosed bool
			err = cupsRepository.DB.QueryRowx(
				`SELECT log_id, cup_index, lad, ink, art, irk, block_number, is_closed FROM maker.cups`).
				Scan(&DBlogId, &DBcupIndex, &DBlad, &DBink, &DBart, &DBirk, &DBblockNumber, &DBisClosed)
			Expect(err).ToNot(HaveOccurred())
			Expect(DBlogId).To(Equal(logId))
			Expect(DBcupIndex).To(Equal(cupIndex))
			Expect(DBlad).To(Equal(lad.Hex()))
			Expect(DBink).To(Equal(ink.Int64()))
			Expect(DBart).To(Equal(art.Int64()))
			Expect(DBirk).To(Equal(irk.Int64()))
			Expect(DBblockNumber).To(Equal(blockNumber))
			Expect(DBisClosed).To(Equal(false))
		})

		It("does not duplicate cups for the same log ID", func() {

			err := logRepository.CreateLogs([]core.Log{{}})
			Expect(err).ToNot(HaveOccurred())

			cupIndex := int64(8675)
			lad := common.StringToAddress("Address")
			ink := big.NewInt(1)
			art := big.NewInt(29259384232352)
			irk := big.NewInt(3)
			cup := fetchers.Cup{
				Lad: lad,
				Ink: ink,
				Art: art,
				Irk: irk,
			}

			var logId int64
			err = logRepository.Get(&logId, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())

			blockNumber := int64(1234)
			err = cupsRepository.CreateCup(logId, cup, blockNumber, false, cupIndex)
			Expect(err).ToNot(HaveOccurred())
			err = cupsRepository.CreateCup(logId, cup, blockNumber, false, cupIndex)
			Expect(err).ToNot(HaveOccurred())

			cups := []DBCup{}
			err = cupsRepository.DB.Select(&cups, `SELECT log_id, cup_index, lad, ink, art, irk, block_number, is_closed FROM maker.cups`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(cups)).To(Equal(1))
		})
	})

	Context("Get cup event", func() {
		var logId int64
		BeforeEach(func() {
			filter := filters.LogFilter{
				Name:      "Cups",
				FromBlock: 0,
				ToBlock:   10,
				Address:   "0x123",
				Topics:    core.Topics{0: "event1=10", 2: "event3=hello"},
			}
			logs := []core.Log{
				{
					BlockNumber: 0,
					TxHash:      "0x1",
					Address:     "0x123",
					Topics:      core.Topics{0: "event1=10", 2: "event3=hello"},
					Index:       0,
					Data:        "",
				},
			}
			err := filterRepository.CreateFilter(filter)
			Expect(err).ToNot(HaveOccurred())

			err = logRepository.CreateLogs(logs)
			logRepository.Get(&logId, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when event has not been seen", func() {
			It("returns unseen events", func() {
				matchingLogs, err := cupsRepository.GetCupEvents()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(matchingLogs)).To(Equal(1))
			})

			It("does not return a seen event", func() {
				cupIndex := int64(8675)
				lad := common.StringToAddress("Address")
				ink := big.NewInt(1)
				art := big.NewInt(29259384232352)
				irk := big.NewInt(3)
				cup := fetchers.Cup{
					Lad: lad,
					Ink: ink,
					Art: art,
					Irk: irk,
				}
				err := cupsRepository.CreateCup(logId, cup, 1234, false, cupIndex)
				Expect(err).ToNot(HaveOccurred())

				matchingLogs, err := cupsRepository.GetCupEvents()
				Expect(err).ToNot(HaveOccurred())
				Expect(len(matchingLogs)).To(Equal(0))
			})
		})
	})
})
