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

package gov_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
	"github.com/8thlight/sai_watcher/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var _ = Describe("Gov transformer", func() {
	It("fetches watched events for 'mold' events", func() {
		watchedEvents := []*core.WatchedEvent{{
			LogID:       0,
			Name:        "",
			BlockNumber: 0,
			Address:     "",
			TxHash:      "",
			Index:       0,
			Topic0:      "",
			Topic1:      "",
			Topic2:      "",
			Topic3:      "",
			Data:        "",
		}}
		mockWatchedEventsRepo := test_helpers.MockWatchedEventsRepository{
			ReturnEvents: watchedEvents,
		}
		transformer := gov.GovTransformer{
			Blockchain:             &test_helpers.MockBlockchain{},
			WatchedEventRepository: &mockWatchedEventsRepo,
			Fetcher:                &test_helpers.MockGovFetcher{},
			GovRepository:          &test_helpers.MockGovRepository{},
		}

		transformer.Execute()

		Expect(len(mockWatchedEventsRepo.EventNames)).To(Equal(1))
		Expect(mockWatchedEventsRepo.EventNames[0]).To(Equal(gov.GovFilter.Name))
	})

	It("fetches gov data for corresponding watched event", func() {
		blockNumber := int64(12345)
		returnEvents := []*core.WatchedEvent{
			{
				LogID:       0,
				Name:        "",
				BlockNumber: blockNumber,
				Address:     "",
				TxHash:      "",
				Index:       0,
				Topic0:      gov.MoldActionHex,
				Topic1:      "",
				Topic2:      "",
				Topic3:      "",
				Data:        "",
			},
		}
		mockWatchedEventsRepo := test_helpers.MockWatchedEventsRepository{ReturnEvents: returnEvents}
		mockFetcher := test_helpers.MockGovFetcher{}
		transformer := gov.GovTransformer{
			Blockchain:             &test_helpers.MockBlockchain{},
			WatchedEventRepository: &mockWatchedEventsRepo,
			Fetcher:                &mockFetcher,
			GovRepository:          &test_helpers.MockGovRepository{},
		}

		transformer.Execute()

		Expect(len(mockFetcher.AxeCalls)).To(Equal(1))
		Expect(mockFetcher.AxeCalls[0]).To(Equal(blockNumber))
		Expect(len(mockFetcher.CapCalls)).To(Equal(1))
		Expect(mockFetcher.CapCalls[0]).To(Equal(blockNumber))
		Expect(len(mockFetcher.FeeCalls)).To(Equal(1))
		Expect(mockFetcher.FeeCalls[0]).To(Equal(blockNumber))
		Expect(len(mockFetcher.GapCalls)).To(Equal(1))
		Expect(mockFetcher.GapCalls[0]).To(Equal(blockNumber))
		Expect(len(mockFetcher.MatCalls)).To(Equal(1))
		Expect(mockFetcher.MatCalls[0]).To(Equal(blockNumber))
		Expect(len(mockFetcher.TaxCalls)).To(Equal(1))
		Expect(mockFetcher.TaxCalls[0]).To(Equal(blockNumber))
	})
})
