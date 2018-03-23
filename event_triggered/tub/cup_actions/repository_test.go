package cup_actions_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

var _ = Describe("Cup Actions Repository", func() {
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
		db.Query(`DELETE FROM maker.cup_action`)
		db.Query(`DELETE FROM logs`)
		logRepository = repositories.LogRepository{DB: db}
	})

	It("persists a cup action model", func() {
		err := logRepository.CreateLogs([]core.Log{{}})
		Expect(err).ToNot(HaveOccurred())
		var logID int64
		err = logRepository.Get(&logID, `Select id from logs`)
		Expect(err).ToNot(HaveOccurred())
		id := int64(12345)
		tx := "Transaction"
		act := "open"
		arg := "Arg"
		lad := "Lad"
		ink := "123"
		art := "456"
		ire := "789"
		guy := "Guy"
		block := int64(54321)
		cupAction := cup_actions.CupActionModel{
			ID:              id,
			TransactionHash: tx,
			Act:             act,
			Arg:             arg,
			Lad:             lad,
			Ink:             ink,
			Art:             art,
			Ire:             ire,
			Guy:             guy,
			Block:           block,
			Deleted:         true,
		}

		err = cup_actions.CupActionsRepository{DB: db}.CreateCupAction(cupAction, logID)
		Expect(err).ToNot(HaveOccurred())

		var DBid int64
		var DBtransactionHash string
		var DBact string
		var DBarg string
		var DBlad string
		var DBink string
		var DBart string
		var DBire string
		var DBblock int64
		var DBdeleted bool
		err = cup_actions.CupActionsRepository{DB: db}.DB.QueryRowx(
			`SELECT id, tx, act, arg, lad, ink, art, ire, block, deleted FROM maker.cup_action`).
			Scan(&DBid, &DBtransactionHash, &DBact, &DBarg, &DBlad, &DBink, &DBart, &DBire, &DBblock, &DBdeleted)
		Expect(err).ToNot(HaveOccurred())
		Expect(DBid).To(Equal(id))
		Expect(DBtransactionHash).To(Equal(tx))
		Expect(DBact).To(Equal(act))
		Expect(DBarg).To(Equal(arg))
		Expect(DBlad).To(Equal(lad))
		Expect(DBink).To(Equal(ink))
		Expect(DBart).To(Equal(art))
		Expect(DBire).To(Equal(ire))
		Expect(DBblock).To(Equal(block))
		Expect(DBdeleted).To(BeTrue())
	})

	It("removes a cup action when corresponding log is removed", func() {
		err := logRepository.CreateLogs([]core.Log{{}})
		Expect(err).ToNot(HaveOccurred())
		var logID int64
		err = logRepository.Get(&logID, `Select id from logs`)
		Expect(err).ToNot(HaveOccurred())

		cupAction := cup_actions.CupActionModel{
			ID:              0,
			TransactionHash: "",
			Act:             "open",
			Arg:             "arg",
			Lad:             "",
			Ink:             "1",
			Art:             "2",
			Ire:             "3",
			Block:           0,
			Deleted:         false,
		}

		//confirm newly created cup action is present
		cup_actions_repo := cup_actions.CupActionsRepository{DB: db}
		err = cup_actions_repo.CreateCupAction(cupAction, logID)
		Expect(err).ToNot(HaveOccurred())
		type dbRow struct {
			LogId int64 `db:"log_id"`
			cup_actions.CupActionModel
		}
		result := &dbRow{}
		err = cup_actions_repo.DB.QueryRowx(
			`SELECT * FROM maker.cup_action WHERE log_id = $1`, logID).StructScan(result)
		Expect(err).ToNot(HaveOccurred())
		Expect(result.LogId).To(Equal(logID))

		//log is removed b/c of reorg
		var logCount int64
		_, err = logRepository.DB.Exec(`DELETE FROM logs WHERE id = $1`, logID)
		Expect(err).ToNot(HaveOccurred())
		err = logRepository.Get(&logCount, `SELECT count(*) FROM logs WHERE id = $1`, logID)
		Expect(err).ToNot(HaveOccurred())
		Expect(logCount).To(BeZero())

		//confirm corresponding cup action is removed
		var cupActionCount int
		err = logRepository.DB.QueryRowx(
			`SELECT count(*) FROM maker.cup_action WHERE log_id = $1`, logID).Scan(&cupActionCount)
		Expect(err).ToNot(HaveOccurred())
		Expect(cupActionCount).To(BeZero())

	})
})
