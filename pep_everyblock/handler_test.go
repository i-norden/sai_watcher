package peps_everyblock_test

import (
	"math/big"

	"github.com/8thlight/sai_watcher/pep_everyblock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockPepsRepository struct {
	values       []string
	blockNumbers []int64
	logIds       []int64
}

func (mpr *MockPepsRepository) CreatePep(value string, blockNumber int64) error {
	mpr.values = append(mpr.values, value)
	mpr.blockNumbers = append(mpr.blockNumbers, blockNumber)
	return nil
}

func (mpr *MockPepsRepository) CheckNewPep() ([]*core.WatchedEvent, error) {
	panic("implement me")
}

type FakePepFetcher struct {
	lastBlock *big.Int
	result    peps_everyblock.PeekResult
}

func (fakePepFetcher FakePepFetcher) GetBlockByNumber(blockNumber int64) core.Block {
	panic("implement me")
}

func (fakePepFetcher FakePepFetcher) GetLogs(contract core.Contract, startingBlockNumber *big.Int, endingBlockNumber *big.Int) ([]core.Log, error) {
	panic("implement me")
}

func (fakePepFetcher FakePepFetcher) Node() core.Node {
	panic("implement me")
}

func (fakePepFetcher FakePepFetcher) FetchContractData(abiJSON string, address string, method string, methodArg interface{}, result interface{}, blockNumber int64) error {
	return nil
}

func (fakePepFetcher FakePepFetcher) LastBlock() *big.Int {
	return fakePepFetcher.lastBlock
}

var filterFirstBlock = big.NewInt(5209657)

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
		pepUpdater := peps_everyblock.NewPepHandler(db, &FakePepFetcher{})
		blockchain := &fakeContractDataFetcher{lastBlock: filterFirstBlock}
		pepsRepository := &MockPepsRepository{}
		pepUpdater = &peps_everyblock.Handler{IPepsRepository: pepsRepository, Blockchain: blockchain}

		pepUpdater.Execute()

		Expect(len(pepsRepository.blockNumbers)).To(Equal(1))
		Expect(pepsRepository.blockNumbers[0]).To(Equal(filterFirstBlock.Int64()))
	})

	It("makes call for every block in filter range", func() {
		lastBlock := filterFirstBlock.Int64() + 24
		blockchain := &fakeContractDataFetcher{lastBlock: big.NewInt(lastBlock)}
		pepUpdater := peps_everyblock.NewPepHandler(db, blockchain)
		pepsRepository := &MockPepsRepository{}
		pepUpdater = &peps_everyblock.Handler{IPepsRepository: pepsRepository, Blockchain: blockchain}

		pepUpdater.Execute()

		Expect(err).ToNot(HaveOccurred())
		Expect(len(blockchain.blocknumbers)).To(Equal(25))

	})
})
