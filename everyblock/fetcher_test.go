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

package everyblock_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/everyblock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

type fakeContractDataFetcher struct {
	abis         []string
	addresses    []string
	methods      []string
	methodArgs   []interface{}
	results      []interface{}
	blocknumbers []int64
	lastBlock    *big.Int
}

func (cdf *fakeContractDataFetcher) GetBlockByNumber(blockNumber int64) core.Block {
	panic("implement me")
}

func (cdf *fakeContractDataFetcher) GetLogs(contract core.Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]core.Log, error) {
	panic("implement me")
}

func (cdf *fakeContractDataFetcher) Node() core.Node {
	panic("implement me")
}

func (cdf *fakeContractDataFetcher) LastBlock() *big.Int {
	return cdf.lastBlock
}

func (cdf *fakeContractDataFetcher) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	cdf.abis = append(cdf.abis, abiJSON)
	cdf.addresses = append(cdf.addresses, address)
	cdf.methods = append(cdf.methods, method)
	cdf.methodArgs = append(cdf.methodArgs, methodArg)
	cdf.results = append(cdf.results, result)
	cdf.blocknumbers = append(cdf.blocknumbers, blockNumber)
	return nil
}

var _ = Describe("Medianizer Data Fetcher", func() {
	var infuraIPC string
	BeforeEach(func() {
		infuraIPC = "https://mainnet.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
	})

	Describe("Getting medianizer attributes", func() {
		It("contract data fetcher with correct args", func() {
			blockchain := &fakeContractDataFetcher{}
			client := everyblock.NewFetcher(blockchain)
			blockNumber := int64(5136253)
			var (
				ret0 = new([32]byte)
				ret1 = new(bool)
			)
			expected := &[]interface{}{
				ret0,
				ret1,
			}

			_, err := client.FetchPepData(nil, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(blockchain.abis)).To(Equal(1))
			abiJSON := everyblock.MedianizerABI
			Expect(err).NotTo(HaveOccurred())
			Expect(blockchain.abis[0]).To(Equal(abiJSON))
			Expect(len(blockchain.addresses)).To(Equal(1))
			Expect(blockchain.addresses[0]).To(Equal("0x99041f808d598b782d5a3e498681c2452a31da08"))
			Expect(len(blockchain.methods)).To(Equal(1))
			Expect(blockchain.methods[0]).To(Equal("peek"))
			Expect(len(blockchain.methodArgs)).To(Equal(1))
			Expect(blockchain.methodArgs[0]).To(BeNil())
			Expect(len(blockchain.results)).To(Equal(1))
			Expect(blockchain.results[0]).To(Equal(expected))
			Expect(len(blockchain.blocknumbers)).To(Equal(1))
			Expect(blockchain.blocknumbers[0]).To(Equal(blockNumber))
			Expect(err).NotTo(HaveOccurred())
		})

	})

	It("makes call to real blockchain", func() {
		blockchain := geth.NewBlockchain(infuraIPC)
		client := everyblock.NewFetcher(blockchain)

		result, err := client.FetchPepData(nil, 5136253)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Value.Hex()).To(Equal("0x0000000000000000000000000000000000000000000000359d858309aa630800"))
		Expect(result.Value.String()).To(Equal("989028058420000000000"))
		Expect(result.OK).To(Equal(true))
	})

	Describe(" matching dai service values", func() {
		var (
			blockNumber int64 = 5237067
			pip               = "703.57"
			pep               = "817.88284690765"
			per               = "1.0020921650678054"
		)

		It("returns the correct converted values for a real pep", func() {
			blockchain := geth.NewBlockchain(infuraIPC)
			client := everyblock.NewFetcher(blockchain)

			result, err := client.FetchPepData(nil, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Wad()).To(Equal(pep))

		})
		It("returns the correct converted values for a real pip", func() {
			blockchain := geth.NewBlockchain(infuraIPC)
			client := everyblock.NewFetcher(blockchain)

			result, err := client.FetchPipData(nil, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Wad()).To(Equal(pip))

		})

		It("returns the correct converted values for a real per", func() {
			blockchain := geth.NewBlockchain(infuraIPC)
			client := everyblock.NewFetcher(blockchain)

			result, err := client.FetchPerData(nil, blockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Ray()).To(Equal(per))

		})
	})

})
