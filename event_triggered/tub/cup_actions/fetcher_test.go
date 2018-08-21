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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"math/big"
)

var _ = Describe("Fetcher", func() {
	It("fetches contract data", func() {
		mockBlockchain := test_helpers.MockBlockchain{}
		fetcher := cup_actions.CupFetcher{BlockChain: &mockBlockchain}
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

	It("fetches cup data at the given block height", func() {
		infuraIPC := "https://mainnet.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"
		rawRpcClient, err := rpc.Dial(infuraIPC)
		Expect(err).NotTo(HaveOccurred())
		rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
		ethClient := ethclient.NewClient(rawRpcClient)
		blockChainClient := client.NewEthClient(ethClient)
		node := node.MakeNode(rpcClient)
		transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
		blockChain := geth.NewBlockChain(blockChainClient, node, transactionConverter)
		realFetcher := cup_actions.NewCupFetcher(blockChain)
		cupID := "0x00000000000000000000000000000000000000000000000000000000000002c6"
		args := common.HexToHash(cupID)
		blockNumber := int64(5257349)
		result, err := realFetcher.FetchCupData(args, blockNumber)

		Expect(err).NotTo(HaveOccurred())
		var ink, art, ire big.Int
		ink.SetString("11536577693755896000", 10)
		art.SetString("1991000000000000000000", 10)
		ire.SetString("1989032743588388759284", 10)
		expectedResult := cup_actions.Cup{
			Lad: common.HexToAddress("0x9b4e28020b94b28f9f09ede87f588e89c283cffd"),
			Ink: &ink,
			Art: &art,
			Ire: &ire,
		}

		Expect(result.Lad).To(Equal(expectedResult.Lad))
		Expect(result.Ink.String()).To(Equal(expectedResult.Ink.String()))
		Expect(result.Art.String()).To(Equal(expectedResult.Art.String()))
		Expect(result.Ire.String()).To(Equal(expectedResult.Ire.String()))
	})
})
