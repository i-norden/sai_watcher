package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/vulcanize/vulcanizedb/pkg/config"
)

var rootCmd = &cobra.Command{
    Use: "sai-watcher",
}

var cfg = config.Database{
	Hostname: "localhost",
	Name:     "vulcanize_public",
	Port:     5432,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
