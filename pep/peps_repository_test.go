package pep

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var _ = Describe("Peps Repository", func() {
	var db *postgres.DB
	var pepsRepository PepsRepository
	var logsRepository repositories.LogRepository
	var filterRepository repositories.FilterRepository

	BeforeEach(func() {
		db = postgres.NewTestDB(core.Node{})
		pepsRepository = PepsRepository{DB: db}
		logsRepository = repositories.LogRepository{DB: db}
		filterRepository = repositories.FilterRepository{DB: db}
	})

	Describe("Creating a new pep record", func() {
		It("inserts new pep with data", func() {
			logsRepository.CreateLogs([]core.Log{
				{},
			})
			var logID int64
			err := logsRepository.Get(&logID, `SELECT id FROM logs`)
			Expect(err).ToNot(HaveOccurred())
			pep := DBPep{
				Value:     "10",
				Timestamp: 123,
			}
			err = pepsRepository.CreatePep(pep.Value, 10, logID)
			Expect(err).ToNot(HaveOccurred())

			var id int64
			var logId int64
			var blockNumber int64
			err = pepsRepository.DB.QueryRow(`SELECT 
				id, log_id, block_number FROM peps`).
				Scan(&id, &logId, &blockNumber)
			Expect(err).ToNot(HaveOccurred())
			Expect(id).ToNot(BeNil())
			Expect(logId).ToNot(BeNil())
			Expect(blockNumber).To(Equal(int64(10)))
		})
	})

	Describe("Get pep event", func() {
		var logIDs []int64

		BeforeEach(func() {
			logIDs = []int64{}
			filter := filters.LogFilter{
				Name:      "Peps",
				FromBlock: 0,
				ToBlock:   10,
				Address:   "0x123",
				Topics:    core.Topics{0: "event1=10", 2: "event3=hello"},
			}
			err := filterRepository.CreateFilter(filter)
			Expect(err).ToNot(HaveOccurred())
			pepEvent := core.Log{
				BlockNumber: 8,
				Address:     "0x123",
				Topics:      core.Topics{0: "event1=10", 2: "event3=hello"},
			}
			pepEvent2 := core.Log{
				BlockNumber: 9,
				Address:     "0x123",
				Topics:      core.Topics{0: "event1=10", 2: "event3=hello"},
			}
			err = logsRepository.CreateLogs([]core.Log{
				pepEvent,
				pepEvent2,
			})
			Expect(err).ToNot(HaveOccurred())
			err = logsRepository.Select(&logIDs, `SELECT id FROM logs`)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when event has not been seen", func() {
			It("returns unseen events", func() {
				watchedEvents, err := pepsRepository.CheckNewPep()

				Expect(err).ToNot(HaveOccurred())
				Expect(len(watchedEvents)).To(Equal(2))
			})
		})

		Context("when event has been seen", func() {
			BeforeEach(func() {
				pep := DBPep{
					Value:     "12",
					Timestamp: 321,
				}
				err := pepsRepository.CreatePep(pep.Value, 8, logIDs[0])
				Expect(err).ToNot(HaveOccurred())
			})

			It("does not return the seen event", func() {
				watchedEvents, err := pepsRepository.CheckNewPep()

				Expect(err).ToNot(HaveOccurred())
				Expect(len(watchedEvents)).To(Equal(1))
				Expect(watchedEvents[0].BlockNumber).To(Equal(int64(9)))
			})

			It("does not return any events if all have been seen", func() {
				pep2 := DBPep{
					Value:     "13",
					Timestamp: 322,
				}
				err := pepsRepository.CreatePep(pep2.Value, 9, logIDs[1])
				Expect(err).ToNot(HaveOccurred())

				watchedEvents, err := pepsRepository.CheckNewPep()

				Expect(err).ToNot(HaveOccurred())
				Expect(len(watchedEvents)).To(Equal(0))
			})
		})
	})
})
