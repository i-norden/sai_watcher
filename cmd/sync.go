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

package cmd

import (
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
	"github.com/vulcanize/vulcanizedb/pkg/history"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs vulcanizedb with local ethereum node",
	Long: `Syncs vulcanizedb with local ethereum node.
vulcanizedb sync --startingBlockNumber 0 --config public.toml

Expects ethereum node to be running and requires a .toml config:

  [database]
  name = "vulcanize_public"
  hostname = "localhost"
  port = 5432

  [client]
  ipcPath = "/Users/user/Library/Ethereum/geth.ipc"
`,
	Run: func(cmd *cobra.Command, args []string) {
		sync()
	},
}

const (
	pollingInterval = 7 * time.Second
)

var startingBlockNumber int

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().IntVarP(&startingBlockNumber, "starting-block-number", "s", 0, "Block number to start syncing from")
}

func backFillAllBlocks(blockchain core.BlockChain, blockRepository datastore.BlockRepository, missingBlocksPopulated chan int, startingBlockNumber int64) {
	go func() {
		missingBlocksPopulated <- history.PopulateMissingBlocks(blockchain, blockRepository, startingBlockNumber)
	}()
}

func sync() {
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	rawRpcClient, err := rpc.Dial(ipc)
	if err != nil {
		log.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	client := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(client)
	blockChain := geth.NewBlockChain(client, node, transactionConverter)
	if blockChain.LastBlock().Int64() == 0 {
		log.Fatal("geth initial: state sync not finished")
	}
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal(err)
	}
	blockRepository := repositories.NewBlockRepository(db)
	validator := history.NewBlockValidator(blockChain, blockRepository, 15)

	missingBlocksPopulated := make(chan int)
	_startingBlockNumber := int64(startingBlockNumber)
	go backFillAllBlocks(blockChain, blockRepository, missingBlocksPopulated, _startingBlockNumber)

	for {
		select {
		case <-ticker.C:
			window := validator.ValidateBlocks()
			validator.Log(os.Stdout, window)
		case <-missingBlocksPopulated:
			go backFillAllBlocks(blockChain, blockRepository, missingBlocksPopulated, _startingBlockNumber)
		}
	}
}
