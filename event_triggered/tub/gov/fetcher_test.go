package gov_test

import (
	"github.com/8thlight/sai_watcher/event_triggered/tub"
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
	"github.com/8thlight/sai_watcher/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

var _ = Describe("Gov Fetcher", func() {

	It("fetches axe", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := gov.GovFetcher{&mockBlockchain}
		blockNumber := int64(12345)

		result, err := fetcher.FetchAxe(blockNumber)

		Expect(mockBlockchain.Methods[0]).To(Equal("axe"))
		AssertFetchCalled(err, blockNumber, &result, mockBlockchain)
	})

	It("fetches cap", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := gov.GovFetcher{&mockBlockchain}
		blockNumber := int64(12345)

		result, err := fetcher.FetchCap(blockNumber)

		Expect(mockBlockchain.Methods[0]).To(Equal("cap"))
		AssertFetchCalled(err, blockNumber, &result, mockBlockchain)
	})

	It("fetches mat", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := gov.GovFetcher{&mockBlockchain}
		blockNumber := int64(12345)

		result, err := fetcher.FetchMat(blockNumber)

		Expect(mockBlockchain.Methods[0]).To(Equal("mat"))
		AssertFetchCalled(err, blockNumber, &result, mockBlockchain)
	})

	It("fetches tax", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := gov.GovFetcher{&mockBlockchain}
		blockNumber := int64(12345)

		result, err := fetcher.FetchTax(blockNumber)

		Expect(mockBlockchain.Methods[0]).To(Equal("tax"))
		AssertFetchCalled(err, blockNumber, &result, mockBlockchain)
	})

	It("fetches fee", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := gov.GovFetcher{&mockBlockchain}
		blockNumber := int64(12345)

		result, err := fetcher.FetchFee(blockNumber)

		Expect(mockBlockchain.Methods[0]).To(Equal("fee"))
		AssertFetchCalled(err, blockNumber, &result, mockBlockchain)
	})

	It("fetches gap", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := gov.GovFetcher{&mockBlockchain}
		blockNumber := int64(12345)

		result, err := fetcher.FetchGap(blockNumber)

		Expect(mockBlockchain.Methods[0]).To(Equal("gap"))
		AssertFetchCalled(err, blockNumber, &result, mockBlockchain)
	})

	Describe("Gov fetcher fetch from real blockchain", func() {
		var (
			infuraIPC         = "https://mainnet.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
			blockNumber int64 = 4755945
			cap_              = "50000000000000000000000000"
			mat               = "1500000000000000000000000000"
			tax               = "1000000000000000000000000000"
			fee               = "1000000000158153903837946257"
			axe               = "1130000000000000000000000000"
			gap               = "1000000000000000000"
		)

		It("returns the correct converted values for a real gov", func() {
			blockchain := geth.NewBlockchain(infuraIPC)
			client := gov.GovFetcher{Blockchain: blockchain}

			result, err := client.FetchAxe(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(axe))
			result, err = client.FetchCap(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(cap_))
			result, err = client.FetchFee(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(fee))
			result, err = client.FetchGap(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(gap))
			result, err = client.FetchMat(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(mat))
			result, err = client.FetchTax(blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.String()).To(Equal(tax))

		})

	})
})

func AssertFetchCalled(err error, blockNumber int64, result interface{}, mockBlockchain test_helpers.MockBlockchain) {
	Expect(err).NotTo(HaveOccurred())
	Expect(len(mockBlockchain.Methods)).To(Equal(1))
	Expect(len(mockBlockchain.MethodArgs)).To(Equal(1))
	Expect(mockBlockchain.MethodArgs[0]).To(BeNil())
	Expect(len(mockBlockchain.Addresses)).To(Equal(1))
	Expect(mockBlockchain.Addresses[0]).To(Equal(tub.TubContractAddress))
	Expect(len(mockBlockchain.AbiJSONs)).To(Equal(1))
	Expect(mockBlockchain.AbiJSONs[0]).To(Equal(tub.TubContractABI))
	Expect(len(mockBlockchain.BlockNumbers)).To(Equal(1))
	Expect(mockBlockchain.BlockNumbers[0]).To(Equal(blockNumber))
	Expect(len(mockBlockchain.Results)).To(Equal(1))
	Expect(mockBlockchain.Results[0]).To(Equal(result))
}
