package main

import (
    "github.com/vulcanize/vulcanizedb/pkg/watchers"
    "github.com/vulcanize/vulcanizedb/pkg/geth"
    "github.com/vulcanize/vulcanizedb/utils"
    "github.com/vulcanize/vulcanizedb/pkg/config"
    "github.com/8thlight/sai_watcher/pep"
)
var cfg = config.Database{
    Hostname: "localhost",
    Name:     "vulcanize_public",
    Port:     5432,
}

func main() {
    blockchain := geth.NewBlockchain("/Users/mattkrump/Library/Ethereum/geth.ipc")
    db := utils.LoadPostgres(cfg, blockchain.Node())
    watcher := watchers.Watcher{
        DB:         db,
        Blockchain: blockchain,
    }
    watcher.AddHandlers(pep.HandlerInitializers())
    watcher.Execute()
}
