package cmd

import (
	"log"

	"time"

	"github.com/8thlight/sai_watcher/event_triggered"
	"github.com/8thlight/sai_watcher/everyblock"
	"github.com/spf13/cobra"
	"github.com/vulcanize/vulcanizedb/libraries/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
)

// getEventsCmd represents the getEvents command
var getEventsCmd = &cobra.Command{
	Use:   "getEvents",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getEvents()
	},
}

func init() {
	rootCmd.AddCommand(getEventsCmd)
}

func getEvents() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	blockchain := geth.NewBlockchain(ipc)
	db, err := postgres.NewDB(databaseConfig, blockchain.Node())
	if err != nil {
		log.Fatal("DB")
	}
	watcher := shared.Watcher{
		DB:         *db,
		Blockchain: blockchain,
	}
	watcher.AddHandlers(event_triggered.HandlerInitializers())
	watcher.AddHandlers(everyblock.HandlerInitializers())
	for range ticker.C {
		watcher.Execute()
	}

}
