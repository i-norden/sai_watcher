package fetchers

import (
	"math/big"

	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/8thlight/sai_watcher/utils"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/ethereum/go-ethereum/common"
)

type fakeCupDataFetcher struct {
	abis         []string
	addresses    []string
	methods      []string
	methodArgs   []interface{}
	results      []interface{}
	blocknumbers []int64
}

func (cdf *fakeCupDataFetcher) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	cdf.abis = append(cdf.abis, abiJSON)
	cdf.addresses = append(cdf.addresses, address)
	cdf.methods = append(cdf.methods, method)
	cdf.methodArgs = append(cdf.methodArgs, methodArg)
	cdf.results = append(cdf.results, result)
	cdf.blocknumbers = append(cdf.blocknumbers, blockNumber)
	return nil
}

var _ = Describe("Cup Data Fetcher", func() {
	Describe("Getting cup attributes", func() {
		It("correctly decodes byte array to JSON", func() {
			blockchain := &fakeCupDataFetcher{}
			client := NewFetcher(blockchain)
			args := common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000002c6")
			blockNumber := int64(5136253)

			_, err := client.FetchCupData(args, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(blockchain.abis)).To(Equal(1))
			expectedAbiPath := filepath.Join(utils.ProjectRoot(), "cup", "saitub.json")
			expectedAbi, err := geth.ReadAbiFile(expectedAbiPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(blockchain.abis[0]).To(Equal(expectedAbi))
			Expect(len(blockchain.addresses)).To(Equal(1))
			Expect(blockchain.addresses[0]).To(Equal("0x448a5065aebb8e423f0896e6c5d525c040f59af3"))
			Expect(len(blockchain.methods)).To(Equal(1))
			Expect(blockchain.methods[0]).To(Equal("cups"))
			Expect(len(blockchain.methodArgs)).To(Equal(1))
			Expect(blockchain.methodArgs[0]).To(Equal(args))
			Expect(len(blockchain.results)).To(Equal(1))
			Expect(blockchain.results[0]).To(Equal(&Cup{}))
			Expect(len(blockchain.blocknumbers)).To(Equal(1))
			Expect(blockchain.blocknumbers[0]).To(Equal(blockNumber))
		})

	})

	It("makes call to passed blockchain", func() {
		config, err := config.NewConfig("infura")
		Expect(err).NotTo(HaveOccurred())
		blockchain := geth.NewBlockchain(config.Client.IPCPath)
		client := NewFetcher(blockchain)
		args := common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000002c6")

		contractData, err := client.FetchCupData(args, 5136253)

		Expect(err).NotTo(HaveOccurred())
		ink := new(big.Int)
		ink.SetString("3825000000000000000", 10)
		Expect(ink).To(Equal(contractData.Ink))
		art := new(big.Int)
		art.SetString("720000000000000000000", 10)
		Expect(art).To(Equal(contractData.Art))
		irk := new(big.Int)
		irk.SetString("719369287647780430799", 10)
		Expect(irk).To(Equal(contractData.Irk))
		Expect("0x9b4e28020B94B28f9f09edE87F588e89c283cFFD").To(Equal(contractData.Lad.Hex()))
	})
})
