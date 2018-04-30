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
	"github.com/8thlight/sai_watcher/event_triggered/tub"
	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	"github.com/8thlight/sai_watcher/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fetcher", func() {
	It("fetches contract data", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := cup_actions.CupFetcher{Blockchain: &mockBlockchain}
		methodArg := "methodArg"
		blockNumber := int64(12345)

		cup, err := fetcher.FetchCupData(methodArg, blockNumber)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(mockBlockchain.AbiJSONs)).To(Equal(1))
		Expect(mockBlockchain.AbiJSONs[0]).To(Equal(tub.TubContractABI))
		Expect(len(mockBlockchain.Addresses)).To(Equal(1))
		Expect(mockBlockchain.Addresses[0]).To(Equal(tub.TubContractAddress))
		Expect(len(mockBlockchain.Methods)).To(Equal(1))
		Expect(mockBlockchain.Methods[0]).To(Equal(cup_actions.CupsContractMethod))
		Expect(len(mockBlockchain.MethodArgs)).To(Equal(1))
		Expect(mockBlockchain.MethodArgs[0]).To(Equal(methodArg))
		Expect(len(mockBlockchain.Results)).To(Equal(1))
		Expect(mockBlockchain.Results[0]).To(Equal(cup))
		Expect(len(mockBlockchain.BlockNumbers)).To(Equal(1))
		Expect(mockBlockchain.BlockNumbers[0]).To(Equal(blockNumber))
		Expect(1).To(Equal(1))
	})
})
