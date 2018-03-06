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
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
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
	blockchain := geth.NewBlockchain(ipc)
	db, err := postgres.NewDB(databaseConfig, blockchain.Node())
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
		fmt.Println(written)
	}

	ipfsCommand := exec.Command("ipfs", "dag", "put", filename)
	output, err := ipfsCommand.Output()
	if err != nil {
		fmt.Println("Unable to write to IPFS", err)
	}
	fmt.Println("Created IPFS hash: ", string(output))
}
