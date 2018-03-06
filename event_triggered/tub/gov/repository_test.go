package gov_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Gov Repository", func() {
	var db *postgres.DB
	var logRepository repositories.LogRepository

	BeforeEach(func() {
		var err error
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())
		db.Query(`DELETE FROM maker.gov`)
		db.Query(`DELETE FROM logs`)
		logRepository = repositories.LogRepository{DB: db}
	})

	Describe("Creating a gov", func() {
		It("persists a gov model", func() {
			err := logRepository.CreateLogs([]core.Log{{}})
			Expect(err).ToNot(HaveOccurred())
			var logID int64
			err = logRepository.Get(&logID, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())
			block := int64(54321)
			tx := "TX"
			var_ := "VAR"
			arg_ := "1"
			guy := "GUY"
			cap_ := "2"
			mat := "3"
			tax := "4"
			fee := "5"
			axe := "6"
			gap := "7"
			govModel := gov.GovModel{
				Block: block,
				Tx:    tx,
				Var:   var_,
				Arg:   arg_,
				Guy:   guy,
				Cap:   cap_,
				Mat:   mat,
				Tax:   tax,
				Fee:   fee,
				Axe:   axe,
				Gap:   gap,
			}
			err = gov.DataStore{DB: db}.CreateGov(&govModel, logID)
			Expect(err).ToNot(HaveOccurred())

			type dbGov struct {
				LogID int `db:"log_id"`
				*gov.GovModel
			}
			result := &dbGov{}
			err = gov.DataStore{DB: db}.DB.QueryRowx(
				`SELECT * FROM maker.gov`).StructScan(result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.LogID).To(Not(BeNil()))
			Expect(result.Tx).To(Equal(tx))
			Expect(result.Var).To(Equal(var_))
			Expect(result.Arg).To(Equal(arg_))
			Expect(result.Guy).To(Equal(guy))
			Expect(result.Cap).To(Equal(cap_))
			Expect(result.Mat).To(Equal(mat))
			Expect(result.Tax).To(Equal(tax))
			Expect(result.Fee).To(Equal(fee))
			Expect(result.Axe).To(Equal(axe))
			Expect(result.Gap).To(Equal(gap))
		})

		It("does not duplicate govs from the same transaction", func() {
			err := logRepository.CreateLogs([]core.Log{{}, {}})
			Expect(err).ToNot(HaveOccurred())
			var logIDs []int64
			err = logRepository.DB.Select(&logIDs, `Select id from logs`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(logIDs)).To(Equal(2))
			govOne := gov.GovModel{
				Block: 0,
				Tx:    "Same",
				Var:   "",
				Arg:   "1",
				Guy:   "",
				Cap:   "1",
				Mat:   "1",
				Tax:   "1",
				Fee:   "1",
				Axe:   "1",
				Gap:   "1",
			}
			govTwo := gov.GovModel{
				Block: 0,
				Tx:    "Same",
				Var:   "",
				Arg:   "2",
				Guy:   "",
				Cap:   "2",
				Mat:   "2",
				Tax:   "2",
				Fee:   "2",
				Axe:   "2",
				Gap:   "2",
			}
			repo := gov.DataStore{db}

			err = repo.CreateGov(&govOne, logIDs[0])
			err = repo.CreateGov(&govTwo, logIDs[1])

			var govs []gov.GovModel
			err = repo.DB.Select(&govs, `SELECT block, tx, var, arg, guy, cap, mat, tax, fee, axe, gap FROM maker.gov`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(govs)).To(Equal(1))
		})
	})

	Describe("Handling reorgs", func() {
		It("removes a gov when corresponding log is removed", func() {
			err := logRepository.CreateLogs([]core.Log{{}})
			Expect(err).ToNot(HaveOccurred())
			var logID int64
			err = logRepository.Get(&logID, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())

			govEvent := &gov.GovModel{
				Block: 0,
				Tx:    "TX",
				Var:   "VAR",
				Arg:   "1",
				Guy:   "GUY",
				Cap:   "2",
				Mat:   "3",
				Tax:   "4",
				Fee:   "5",
				Axe:   "6",
				Gap:   "7",
			}

			//confirm newly created gov is present
			gov_repo := gov.DataStore{DB: db}
			err = gov_repo.CreateGov(govEvent, logID)
			Expect(err).ToNot(HaveOccurred())
			type dbRow struct {
				LogId int64 `db:"log_id"`
				gov.GovModel
			}
			result := &dbRow{}
			err = gov_repo.DB.QueryRowx(
				`SELECT * FROM maker.gov WHERE log_id = $1`, logID).StructScan(result)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.LogId).To(Equal(logID))

			//log is removed b/c of reorg
			var logCount int64
			_, err = logRepository.DB.Exec(`DELETE FROM logs WHERE id = $1`, logID)
			Expect(err).ToNot(HaveOccurred())
			err = logRepository.Get(&logCount, `SELECT count(*) FROM logs WHERE id = $1`, logID)
			Expect(err).ToNot(HaveOccurred())
			Expect(logCount).To(BeZero())

			//confirm corresponding gov is removed
			var govCount int
			err = logRepository.DB.QueryRowx(
				`SELECT count(*) FROM maker.gov WHERE log_id = $1`, logID).Scan(&govCount)
			Expect(err).ToNot(HaveOccurred())
			Expect(govCount).To(BeZero())
		})
	})

	Describe("Fetching all govs", func() {
		It("returns all govs", func() {
			err := logRepository.CreateLogs([]core.Log{{}, {}})
			Expect(err).ToNot(HaveOccurred())
			var logIDs []int64
			err = logRepository.DB.Select(&logIDs, `Select id from logs`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(logIDs)).To(Equal(2))

			gov_repo := gov.DataStore{db}
			blockNumberOne := int64(1)
			govOne := gov.GovModel{
				Block: blockNumberOne,
				Tx:    "TX1",
				Var:   "VarOne",
				Arg:   "1",
				Guy:   "GuyOne",
				Cap:   "2",
				Mat:   "3",
				Tax:   "4",
				Fee:   "5",
				Axe:   "6",
				Gap:   "7",
			}
			err = gov_repo.CreateGov(&govOne, logIDs[0])
			Expect(err).NotTo(HaveOccurred())
			blockNumberTwo := int64(2)
			govTwo := gov.GovModel{
				Block: blockNumberTwo,
				Tx:    "TX2",
				Var:   "VarTwo",
				Arg:   "8",
				Guy:   "GuyTwo",
				Cap:   "9",
				Mat:   "10",
				Tax:   "11",
				Fee:   "12",
				Axe:   "13",
				Gap:   "14",
			}
			err = gov_repo.CreateGov(&govTwo, logIDs[1])
			Expect(err).NotTo(HaveOccurred())

			results, err := gov_repo.GetAllGovData()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(2))
			Expect(results[0]).To(Equal(govOne))
			Expect(results[1]).To(Equal(govTwo))
		})
	})
})
