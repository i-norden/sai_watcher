package peps_everyblock_test

import (
	"path/filepath"

	"math/big"

	"github.com/8thlight/sai_watcher/pep_everyblock"
	"github.com/8thlight/sai_watcher/utils"
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
			client := peps_everyblock.NewFetcher(blockchain)
			blockNumber := int64(5136253)
			var (
				ret0 = new([32]byte)
				ret1 = new(bool)
			)
			expected := &[]interface{}{
				ret0,
				ret1,
			}

			_, err := client.FetchContractData(nil, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(blockchain.abis)).To(Equal(1))
			abiJSON, err := geth.ReadAbiFile(filepath.Join(utils.ProjectRoot(), "pep_everyblock", "medianizer.json"))
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
		client := peps_everyblock.NewFetcher(blockchain)

		result, err := client.FetchContractData(nil, 5136253)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Value.Hex()).To(Equal("0x0000000000000000000000000000000000000000000000359d858309aa630800"))
		Expect(result.Value.String()).To(Equal("989028058420000000000"))
		Expect(result.OK).To(Equal(true))
	})

})
