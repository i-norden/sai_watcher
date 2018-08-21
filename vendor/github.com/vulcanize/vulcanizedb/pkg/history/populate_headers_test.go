package history_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

var _ = Describe("Populating headers", func() {

	var headerRepository *fakes.MockHeaderRepository

	BeforeEach(func() {
		headerRepository = fakes.NewMockHeaderRepository()
	})

	Describe("When 1 missing header", func() {

		It("returns number of headers added", func() {
			blockChain := fakes.NewMockBlockChain()
			blockChain.SetLastBlock(big.NewInt(2))
			headerRepository.SetMissingBlockNumbers([]int64{2})

			headersAdded := history.PopulateMissingHeaders(blockChain, headerRepository, 1)

			Expect(headersAdded).To(Equal(1))
		})
	})

	It("adds missing headers to the db", func() {
		blockChain := fakes.NewMockBlockChain()
		blockChain.SetLastBlock(big.NewInt(2))
		headerRepository.SetMissingBlockNumbers([]int64{2})

		history.PopulateMissingHeaders(blockChain, headerRepository, 1)

		headerRepository.AssertCreateOrUpdateHeaderCallCountAndPassedBlockNumbers(1, []int64{2})
	})
})
