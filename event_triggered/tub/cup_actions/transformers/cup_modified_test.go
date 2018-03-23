package transformers_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions/transformers"
	"github.com/8thlight/sai_watcher/test_helpers"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("Cup Modified transformer", func() {
	It("fetches watched events for cup modified log filters", func() {
		mockEventsRepo := test_helpers.MockWatchedEventsRepository{}
		transformer := transformers.CupModifiedTransformer{
			Blockchain:             nil,
			WatchedEventRepository: &mockEventsRepo,
			CupActionsRepository:   &test_helpers.MockCupActionsRepository{},
			Fetcher:                &test_helpers.MockCupFetcher{},
		}

		transformer.Execute()

		Expect(len(mockEventsRepo.EventNames)).To(Equal(7))
		Expect(mockEventsRepo.EventNames[0]).To(Equal(transformers.CupModifiedFilters[0].Name))
		Expect(mockEventsRepo.EventNames[1]).To(Equal(transformers.CupModifiedFilters[1].Name))
		Expect(mockEventsRepo.EventNames[2]).To(Equal(transformers.CupModifiedFilters[2].Name))
		Expect(mockEventsRepo.EventNames[3]).To(Equal(transformers.CupModifiedFilters[3].Name))
		Expect(mockEventsRepo.EventNames[4]).To(Equal(transformers.CupModifiedFilters[4].Name))
		Expect(mockEventsRepo.EventNames[5]).To(Equal(transformers.CupModifiedFilters[5].Name))
		Expect(mockEventsRepo.EventNames[6]).To(Equal(transformers.CupModifiedFilters[6].Name))
	})

	It("fetches cup data for corresponding watched event", func() {
		blockNumber := int64(12345)
		returnEvents := []*core.WatchedEvent{
			{
				LogID:       0,
				Name:        "",
				BlockNumber: blockNumber,
				Address:     "",
				TxHash:      "",
				Index:       0,
				Topic0:      "",
				Topic1:      "",
				Topic2:      "",
				Topic3:      "0x000000000000000000000000000000000000000000000000172d0826e6e70000",
				Data:        "",
			},
		}
		mockWatchedEventsRepo := test_helpers.MockWatchedEventsRepository{ReturnEvents: returnEvents}
		mockFetcher := test_helpers.MockCupFetcher{}
		transformer := transformers.CupModifiedTransformer{
			CupActionsRepository:   &test_helpers.MockCupActionsRepository{},
			WatchedEventRepository: &mockWatchedEventsRepo,
			FilterRepository:       nil,
			Fetcher:                &mockFetcher,
			Blockchain:             nil,
		}

		transformer.Execute()

		Expect(len(mockFetcher.BlockNumbers)).To(Equal(1))
		Expect(mockFetcher.BlockNumbers[0]).To(Equal(blockNumber))
	})

	It("persists cup action with fetched cup data", func() {
		logID := int64(12345)
		cupID := "0x00000000000000000000000000000000000000000000000000000000000002c6"
		transactionHash := "0xd88ff63af55cb4c206bd2f4dc422552cb24c2f0dc5803c932fc655951a80dc3a"
		blockNumber := int64(5257349)
		returnEvents := []*core.WatchedEvent{
			{
				LogID:       logID,
				Name:        "",
				BlockNumber: blockNumber,
				Address:     transformers.TubContractAddress,
				TxHash:      transactionHash,
				Index:       0,
				Topic0:      transformers.LockActionHex,
				Topic1:      "0x0000000000000000000000009b4e28020b94b28f9f09ede87f588e89c283cffd",
				Topic2:      cupID,
				Topic3:      "0x00000000000000000000000000000000000000000000000007a1fe1602770000",
				Data:        "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044b3b77a5100000000000000000000000000000000000000000000000000000000000002c600000000000000000000000000000000000000000000000007a1fe1602770000",
			},
		}
		mockWatchedEventsRepo := test_helpers.MockWatchedEventsRepository{ReturnEvents: returnEvents}
		ink := big.NewInt(0)
		ink.SetString("11536577693755896000", 10)
		art := big.NewInt(0)
		art.SetString("1991000000000000000000", 10)
		ire := big.NewInt(0)
		ire.SetString("1989032743588388759284", 10)
		returnVal := cup_actions.Cup{
			Lad: common.HexToAddress(returnEvents[0].Topic1),
			Ink: ink,
			Art: art,
			Ire: ire,
		}
		mockFetcher := test_helpers.MockCupFetcher{ReturnVal: returnVal}
		mockCupActionsRepository := test_helpers.MockCupActionsRepository{}
		transformer := transformers.CupModifiedTransformer{
			CupActionsRepository:   &mockCupActionsRepository,
			WatchedEventRepository: &mockWatchedEventsRepo,
			FilterRepository:       nil,
			Fetcher:                &mockFetcher,
			Blockchain:             nil,
		}

		transformer.Execute()

		Expect(len(mockCupActionsRepository.CupActions)).To(Equal(1))
		cupAction := mockCupActionsRepository.CupActions[0]
		Expect(cupAction.ID).To(Equal(shared.HexToInt64(cupID)))
		Expect(cupAction.TransactionHash).To(Equal(transactionHash))
		Expect(cupAction.Act).To(Equal("lock"))
		Expect(cupAction.Arg).To(Equal("0.55"))
		Expect(cupAction.Lad).To(Equal("0x9b4e28020b94b28f9f09ede87f588e89c283cffd"))
		Expect(cupAction.Ink).To(Equal("11.536577693755897"))
		Expect(cupAction.Art).To(Equal("1991"))
		Expect(cupAction.Ire).To(Equal("1989.0327435883887"))
		Expect(cupAction.Guy).To(Equal("0x9b4e28020b94b28f9f09ede87f588e89c283cffd"))
		Expect(cupAction.Block).To(Equal(blockNumber))
		Expect(cupAction.Deleted).To(BeFalse())
		Expect(len(mockCupActionsRepository.LogIDs)).To(Equal(1))
		Expect(mockCupActionsRepository.LogIDs[0]).To(Equal(logID))
	})
})
