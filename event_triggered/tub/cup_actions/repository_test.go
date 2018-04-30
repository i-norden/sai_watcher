// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cup_actions_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/models"
	"github.com/8thlight/sai_watcher/everyblock"
	"github.com/8thlight/sai_watcher/test_helpers"
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
		db.Query(`DELETE FROM logs`)
		db.Query(`DELETE FROM maker.cup_action`)
		db.Query(`DELETE FROM maker.peps_everyblock`)
		db.Query(`DELETE FROM public.blocks`)
		logRepository = repositories.LogRepository{DB: db}
	})

	Describe("Creating a cup action", func() {
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
			cupAction := models.CupAction{
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

		It("does not allow more than one cup action (with same act, arg) per transaction", func() {
			err := logRepository.CreateLogs([]core.Log{{}, {}, {}})
			Expect(err).ToNot(HaveOccurred())
			var logIDs []int64
			err = logRepository.Select(&logIDs, `Select id from logs`)

			id := int64(12345)
			tx := "Transaction"
			lad := "Lad"
			ink := "123"
			art := "456"
			ire := "789"
			guy := "Guy"
			block := int64(54321)
			cupActionUniqueOne := models.CupAction{
				ID:              id,
				TransactionHash: tx,
				Act:             "open",
				Arg:             "arg1",
				Lad:             lad,
				Ink:             ink,
				Art:             art,
				Ire:             ire,
				Guy:             guy,
				Block:           block,
				Deleted:         true,
			}
			cupActionUniqueTwo := models.CupAction{
				ID:              id,
				TransactionHash: tx,
				Act:             "open",
				Arg:             "arg2",
				Lad:             lad,
				Ink:             ink,
				Art:             art,
				Ire:             ire,
				Guy:             guy,
				Block:           block,
				Deleted:         true,
			}
			cupActionDuplicate := cupActionUniqueOne

			cup_actions_repo := cup_actions.CupActionsRepository{DB: db}
			err = cup_actions_repo.CreateCupAction(cupActionUniqueOne, logIDs[0])
			Expect(err).ToNot(HaveOccurred())
			err = cup_actions_repo.CreateCupAction(cupActionUniqueTwo, logIDs[1])
			Expect(err).ToNot(HaveOccurred())
			err = cup_actions_repo.CreateCupAction(cupActionDuplicate, logIDs[2])
			Expect(err).ToNot(HaveOccurred())

			var cupCount int
			cup_actions_repo.DB.Get(&cupCount, `Select count(*) from maker.cup_action`)
			Expect(cupCount).To(Equal(2))
		})
	})

	Describe("Handling reorgs", func() {
		It("removes a cup action when corresponding log is removed", func() {
			err := logRepository.CreateLogs([]core.Log{{}})
			Expect(err).ToNot(HaveOccurred())
			var logID int64
			err = logRepository.Get(&logID, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())

			cupAction := models.CupAction{
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
				models.CupAction
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

	Describe("Getting all cups", func() {
		It("returns all existing cups", func() {
			// Create blocks for cup view
			blocksRepo := repositories.BlockRepository{db}
			blockNumberOne := int64(1234)
			blockOne := core.Block{Number: blockNumberOne}
			err := blocksRepo.CreateOrUpdateBlock(blockOne)
			Expect(err).NotTo(HaveOccurred())
			blockNumberTwo := int64(5678)
			blockTwo := core.Block{Number: blockNumberTwo}
			err = blocksRepo.CreateOrUpdateBlock(blockTwo)
			Expect(err).NotTo(HaveOccurred())

			// Create peps everyblock data for cups view
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
			pepsRepository := everyblock.DataStore{db}
			err = pepsRepository.Create(blockNumberOne, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())
			err = pepsRepository.Create(blockNumberTwo, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())

			// Create logs for cup actions
			err = logRepository.CreateLogs([]core.Log{{}, {}})
			Expect(err).ToNot(HaveOccurred())
			var logIDs []int64
			err = logRepository.DB.Select(&logIDs, `Select id from logs`)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(logIDs)).To(Equal(2))

			// Create cup actions
			cup_actions_repo := cup_actions.CupActionsRepository{db}
			actOne := "wipe"
			artOne := "1"
			idOne := int64(2)
			inkOne := "3"
			ireOne := "4"
			ladOne := "0xlad1"
			cupActionOne := models.CupAction{
				ID:              idOne,
				TransactionHash: "0xtxhash1",
				Act:             actOne,
				Arg:             "",
				Lad:             ladOne,
				Ink:             inkOne,
				Art:             artOne,
				Ire:             ireOne,
				Block:           blockNumberOne,
				Deleted:         false,
				Guy:             "guy1",
			}
			err = cup_actions_repo.CreateCupAction(cupActionOne, logIDs[0])
			Expect(err).NotTo(HaveOccurred())
			actTwo := "lock"
			artTwo := "5"
			idTwo := int64(6)
			inkTwo := "7"
			ireTwo := "8"
			ladTwo := "0xlad2"
			cupActionTwo := models.CupAction{
				ID:              idTwo,
				TransactionHash: "0xtxhash2",
				Act:             actTwo,
				Arg:             "",
				Lad:             ladTwo,
				Ink:             inkTwo,
				Art:             artTwo,
				Ire:             ireTwo,
				Block:           blockNumberTwo,
				Deleted:         true,
				Guy:             "guy2",
			}
			err = cup_actions_repo.CreateCupAction(cupActionTwo, logIDs[1])
			Expect(err).NotTo(HaveOccurred())

			results, err := cup_actions_repo.GetAllCupData()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(2))

			cupOne := results[1]
			Expect(cupOne.Act).To(Equal(actOne))
			Expect(cupOne.Art).To(Equal(artOne))
			Expect(cupOne.Block).To(Equal(blockNumberOne))
			Expect(cupOne.Deleted).To(BeFalse())
			Expect(cupOne.Id).To(Equal(idOne))
			Expect(cupOne.Ink).To(Equal(inkOne))
			Expect(cupOne.Ire).To(Equal(ireOne))
			Expect(cupOne.Lad).To(Equal(ladOne))
			actualCupOnePip := test_helpers.GetFloat(cupOne.Pip)
			Expect(actualCupOnePip.String()).To(Equal(pip.Wad()))
			Expect(cupOne.Per).To(Equal(per.Ray()))
			actualCupOneRatio := test_helpers.GetFloat(*cupOne.Ratio)
			expectedCupOneRatio := getRatio(cupOne, pip, per)
			Expect(actualCupOneRatio.String()).To(Equal(expectedCupOneRatio.String()))
			actualCupOneTab := test_helpers.GetFloat(cupOne.Tab)
			expectedCupOneTab := getTab(cupOne, pip, per)
			Expect(actualCupOneTab.String()).To(Equal(expectedCupOneTab.String()))

			cupTwo := results[0]
			Expect(cupTwo.Act).To(Equal(actTwo))
			Expect(cupTwo.Art).To(Equal(artTwo))
			Expect(cupTwo.Block).To(Equal(blockNumberTwo))
			Expect(cupTwo.Deleted).To(BeTrue())
			Expect(cupTwo.Id).To(Equal(idTwo))
			Expect(cupTwo.Ink).To(Equal(inkTwo))
			Expect(cupTwo.Ire).To(Equal(ireTwo))
			Expect(cupTwo.Lad).To(Equal(ladTwo))
			actualCupTwoPip := test_helpers.GetFloat(cupTwo.Pip)
			Expect(actualCupTwoPip.String()).To(Equal(pip.Wad()))
			Expect(cupTwo.Per).To(Equal(per.Ray()))
			actualCupTwoRatio := test_helpers.GetFloat(*cupTwo.Ratio)
			expectedCupTwoRatio := getRatio(cupTwo, pip, per)
			Expect(actualCupTwoRatio.String()).To(Equal(expectedCupTwoRatio.String()))
			actualCupTwoTab := test_helpers.GetFloat(cupTwo.Tab)
			expectedCupTwoTab := getTab(cupTwo, pip, per)
			Expect(actualCupTwoTab.String()).To(Equal(expectedCupTwoTab.String()))
		})

		It("can scan ratios that are null", func() {
			blocksRepo := repositories.BlockRepository{db}
			blockNumber := int64(1234)
			blockOne := core.Block{Number: blockNumber}
			err := blocksRepo.CreateOrUpdateBlock(blockOne)
			Expect(err).NotTo(HaveOccurred())

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
			pepsRepository := everyblock.DataStore{db}
			err = pepsRepository.Create(blockNumber, pep, pip, per)
			Expect(err).ToNot(HaveOccurred())

			err = logRepository.CreateLogs([]core.Log{{}})
			Expect(err).ToNot(HaveOccurred())
			var logID int64
			err = logRepository.Get(&logID, `Select id from logs`)
			Expect(err).ToNot(HaveOccurred())

			cup_actions_repo := cup_actions.CupActionsRepository{db}
			actOne := "wipe"
			artOne := "0"
			idOne := int64(2)
			inkOne := "0"
			ireOne := "0"
			ladOne := "0xlad1"
			cupActionOne := models.CupAction{
				ID:              idOne,
				TransactionHash: "0xtxhash1",
				Act:             actOne,
				Arg:             "",
				Lad:             ladOne,
				Ink:             inkOne,
				Art:             artOne,
				Ire:             ireOne,
				Block:           blockNumber,
				Deleted:         false,
				Guy:             "guy1",
			}
			err = cup_actions_repo.CreateCupAction(cupActionOne, logID)
			Expect(err).NotTo(HaveOccurred())

			results, err := cup_actions_repo.GetAllCupData()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(Equal(1))
			Expect(results[0].Ratio).To(BeNil())
		})
	})
})

func getRatio(cup models.Cup, pip everyblock.Peek, per everyblock.Per) *big.Float {
	tab := getTab(cup, pip, per)
	artFloat := big.NewFloat(0)
	artFloat.SetString(cup.Art)
	rawRatio := big.NewFloat(0).Quo(tab, artFloat)
	result := big.NewFloat(0).Mul(rawRatio, big.NewFloat(100))
	return result
}

func getTab(cup models.Cup, pip everyblock.Peek, per everyblock.Per) *big.Float {
	expectedPip := big.NewFloat(0)
	expectedPip.SetString(pip.Wad())
	expectedPer := big.NewFloat(0)
	expectedPer.SetString(per.Ray())
	pipTimesPer := big.NewFloat(0).Mul(expectedPip, expectedPer)
	inkFloat := big.NewFloat(0)
	inkFloat.SetString(cup.Ink)
	tab := big.NewFloat(0).Mul(pipTimesPer, inkFloat)
	return tab
}
