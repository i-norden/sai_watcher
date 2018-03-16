# Sai Contract Watcher

## Description
Library code for watching the Sai contract with VulcanizeDB

## Setup
Setup Postgres and Geth - see [VulcanizeDB README](https://github.com/vulcanize/VulcanizeDB/blob/master/README.md)

`./sai_watcher sync --config environments/public.toml --starting-block-number n` (where n is a recent enough block to fetch contract data with your node)

`./sai_watcher getEvents --config environments/public.toml`

`./sai_watcher graphql --config environments/public.toml`

## Run tests
`ginkgo -r`

## Note
Currently this project references a private repo branch of VulcanizeDB, and won't work if you don't have access