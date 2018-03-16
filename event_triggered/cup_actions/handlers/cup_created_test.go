package handlers_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/cup_actions/handlers"
	"github.com/8thlight/sai_watcher/event_triggered/cup_actions/handlers/test_helpers"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var logID = int64(54321)
var blockNumber = int64(12345)
var cupIndexHex = "0x0000000000000000000000000000000000000000000000000000000000000001"
var cupLadHex = "0x000000000000000000000000b0e83c2d71a991017e0116d58c5765abc57384af"
var transactionHash = "0xcf71a7135b820b1436a3a046ca723b4d578a5c7762bc9523f9cb22804c3a75dc"

var returnEvents = []*core.WatchedEvent{
	{
		LogID:       logID,
		Name:        "",
		BlockNumber: blockNumber,
		Address:     "",
		TxHash:      transactionHash,
		Index:       0,
		Topic0:      "",
		Topic1:      cupLadHex,
		Topic2:      "",
		Topic3:      "",
		Data:        cupIndexHex,
	},
}

var _ = Describe("Cup Created Handler", func() {
	It("fetches watched events for cup created filter", func() {
		mockEventsRepo := test_helpers.MockWatchedEventsRepository{ReturnEvents: returnEvents}
		transformer := handlers.CupCreatedHandler{
			Blockchain:             nil,
			WatchedEventRepository: &mockEventsRepo,
			CupActionsRepository:   &test_helpers.MockCupActionsRepository{},
		}

		transformer.Execute()

		Expect(len(mockEventsRepo.EventNames)).To(Equal(1))
		Expect(mockEventsRepo.EventNames[0]).To(Equal(handlers.CupCreatedFilter.Name))
	})

	It("persists cup action for matching event", func() {
		mockEventsRepo := test_helpers.MockWatchedEventsRepository{ReturnEvents: returnEvents}
		mockCupActionsRepo := test_helpers.MockCupActionsRepository{}
		transformer := handlers.CupCreatedHandler{
			Blockchain:             nil,
			WatchedEventRepository: &mockEventsRepo,
			CupActionsRepository:   &mockCupActionsRepo,
		}

		transformer.Execute()

		Expect(len(mockCupActionsRepo.CupActions)).To(Equal(1))
		createdCup := mockCupActionsRepo.CupActions[0]
		Expect(createdCup.Act).To(Equal("open"))
		Expect(createdCup.Arg).To(Equal(""))
		Expect(createdCup.Art).To(Equal("0"))
		Expect(createdCup.Block).To(Equal(blockNumber))
		Expect(createdCup.Deleted).To(BeFalse())
		Expect(createdCup.ID).To(Equal(shared.HexToInt64(cupIndexHex)))
		Expect(createdCup.Ink).To(Equal("0"))
		Expect(createdCup.Ire).To(Equal("0"))
		Expect(createdCup.Lad).To(Equal(common.HexToAddress(cupLadHex).Hex()))
		Expect(createdCup.TransactionHash).To(Equal(transactionHash))
		Expect(len(mockCupActionsRepo.LogIDs)).To(Equal(1))
		Expect(mockCupActionsRepo.LogIDs[0]).To(Equal(logID))
	})
})
