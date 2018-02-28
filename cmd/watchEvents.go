package cmd

import (
"github.com/spf13/cobra"
"github.com/vulcanize/vulcanizedb/pkg/geth"
"github.com/vulcanize/vulcanizedb/pkg/watchers"
"github.com/vulcanize/vulcanizedb/utils"
    "github.com/8thlight/sai_watcher/pep"
    "github.com/8thlight/sai_watcher/cup"
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
    blockchain := geth.NewBlockchain("/Users/mattkrump/Library/Ethereum/geth.ipc")
    db := utils.LoadPostgres(cfg, blockchain.Node())
    watcher := watchers.Watcher{
        DB:         db,
        Blockchain: blockchain,
    }
    watcher.AddHandlers(pep.HandlerInitializers())
    watcher.AddHandlers(cup.HandlerInitializers())
    watcher.Execute()
}