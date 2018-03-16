package cup_actions_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/cup_actions"
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
		err = cup_actions.CupActionsRepository{db}.DB.QueryRowx(
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
})
