package everyblock_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/everyblock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockEveryBlockDataStore struct {
	peps         []everyblock.Peek
	pips         []everyblock.Peek
	pers         []everyblock.Per
	blockNumbers []int64
}

func (mpr *MockEveryBlockDataStore) Get(blockNumber int64) (*everyblock.Row, error) {
	panic("implement me")
}

func (mpr *MockEveryBlockDataStore) Create(blockNumber int64, pep everyblock.Peek, pip everyblock.Peek, per everyblock.Per) error {
	mpr.blockNumbers = append(mpr.blockNumbers, blockNumber)
	mpr.peps = append(mpr.peps, pep)
	mpr.pips = append(mpr.pips, pip)
	mpr.pers = append(mpr.pers, per)
	return nil
}

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

var filterFirstBlock = big.NewInt(everyblock.PepsFilters[0].FromBlock)

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
		db.Query(`DELETE FROM peps_everyblock`)
		db.Query(`DELETE FROM log_filters`)

	})

	It("retrieves a pep for a single block", func() {
		pepUpdater := everyblock.NewPepHandler(db, &FakeBlockchain{})
		blockchain := &fakeContractDataFetcher{lastBlock: filterFirstBlock}
		pepsRepository := &MockEveryBlockDataStore{}
		pepUpdater = &everyblock.Handler{
			Repository: pepsRepository,
			Blockchain: blockchain,
		}

		pepUpdater.Execute()

		Expect(len(pepsRepository.blockNumbers)).To(Equal(1))
		Expect(pepsRepository.blockNumbers[0]).To(Equal(filterFirstBlock.Int64()))
	})

	It("makes call for every block in filter range", func() {
		lastBlock := filterFirstBlock.Int64() + 24
		blockchain := &fakeContractDataFetcher{lastBlock: big.NewInt(lastBlock)}
		pepUpdater := everyblock.NewPepHandler(db, blockchain)
		pepsRepository := &MockEveryBlockDataStore{}
		pepUpdater = &everyblock.Handler{Repository: pepsRepository, Blockchain: blockchain}

		pepUpdater.Execute()

		Expect(err).ToNot(HaveOccurred())
		// 25 * 3 handler calls per block
		Expect(len(blockchain.blocknumbers)).To(Equal(75))

	})
})
