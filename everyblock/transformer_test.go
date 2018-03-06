package everyblock_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/everyblock"
	"github.com/8thlight/sai_watcher/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type FakeBlockchain struct {
	lastBlock *big.Int
	result    everyblock.Peek
}

func (fakePepFetcher FakeBlockchain) GetBlockByNumber(blockNumber int64) core.Block {
	panic("implement me")
}

func (fakePepFetcher FakeBlockchain) GetLogs(contract core.Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]core.Log, error) {
	panic("implement me")
}

func (fakePepFetcher FakeBlockchain) Node() core.Node {
	panic("implement me")
}

func (fakePepFetcher FakeBlockchain) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	return nil
}

func (fakePepFetcher FakeBlockchain) LastBlock() *big.Int {
	return fakePepFetcher.lastBlock
}

var filterFirstBlock = big.NewInt(everyblock.PepsFilter.FromBlock)

var _ bool = Describe("pep updater", func() {
	var db *postgres.DB
	var err error

	BeforeEach(func() {
		db, err = postgres.NewDB(config.Database{
			Hostname: "localhost",
			Name:     "vulcanize_private",
			Port:     5432,
		}, core.Node{})
		Expect(err).NotTo(HaveOccurred())
		db.Query(`DELETE FROM maker.peps_everyblock`)
		db.Query(`DELETE FROM log_filters`)

	})

	It("retrieves a pep for a single block", func() {
		pepUpdater := everyblock.NewPepTransformer(db, &FakeBlockchain{})
		blockchain := &fakeContractDataFetcher{lastBlock: filterFirstBlock}
		int64s := []int64{filterFirstBlock.Int64()}
		pepsRepository := &test_helpers.MockEveryBlockDataStore{MissingBlocksData: int64s}
		pepUpdater = &everyblock.Transformer{
			Repository: pepsRepository,
			Blockchain: blockchain,
		}

		pepUpdater.Execute()

		Expect(len(pepsRepository.BlockNumbers)).To(Equal(1))
		Expect(pepsRepository.BlockNumbers[0]).To(Equal(filterFirstBlock.Int64()))
	})

	It("makes call for every missing block in range", func() {
		lastBlock := filterFirstBlock.Int64() + 5
		blockchain := &fakeContractDataFetcher{lastBlock: big.NewInt(lastBlock)}
		pepUpdater := everyblock.NewPepTransformer(db, blockchain)
		firstBlock := filterFirstBlock.Int64()
		int64s := []int64{firstBlock, firstBlock + 1, firstBlock + 2, firstBlock + 3, firstBlock + 4}
		pepsRepository := &test_helpers.MockEveryBlockDataStore{MissingBlocksData: int64s}
		pepUpdater = &everyblock.Transformer{Repository: pepsRepository, Blockchain: blockchain}

		pepUpdater.Execute()

		Expect(err).ToNot(HaveOccurred())
		// 5 * 3 Transformer calls per block
		Expect(len(blockchain.blocknumbers)).To(Equal(15))

	})
})
