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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/8thlight/sai_watcher/event_triggered/tub/cup_actions"
	"github.com/8thlight/sai_watcher/event_triggered/tub/gov"
	"github.com/8thlight/sai_watcher/everyblock"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

// generateReportCmd represents the generateReport command
var generateReportCmd = &cobra.Command{
	Use:   "generateReport",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateReport()
	},
}

func init() {
	rootCmd.AddCommand(generateReportCmd)
}

func generateReport() {
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
	db, err := postgres.NewDB(databaseConfig, blockChain.Node())
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	cr := cup_actions.CupActionsRepository{db}
	gr := gov.DataStore{db}
	per := everyblock.DataStore{db}
	cups, err := cr.GetAllCupData()
	if err != nil {
		log.Fatal("Could not read cups data", err)
	}
	govs, err := gr.GetAllGovData()
	if err != nil {
		log.Fatal("Could not read gov data", err)
	}
	blocks, err := per.GetAllRows()
	if err != nil {
		log.Fatal("Could not read everyblock data", err)
	}

	results := make(map[string]interface{})
	results["cups"] = cups
	results["govs"] = govs
	results["blocks"] = blocks

	marshalled, err := json.Marshal(results)
	filename := "reports/report_" + time.Now().Format("2018-03-06T15:04:03") + ".txt"

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	written := ioutil.WriteFile(filename, marshalled, 0644)
	if written != nil {
		log.Fatal(written)
	}

	ipfsCommand := exec.Command("ipfs", "dag", "put", filename)
	output, err := ipfsCommand.Output()
	if err != nil {
		log.Fatal("Unable to write to IPFS", err)
	}
	fmt.Println("Created IPFS hash: ", string(output))
}
