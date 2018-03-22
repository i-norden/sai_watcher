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
