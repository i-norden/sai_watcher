# Sai Contract Watcher

## Description
Library code for watching the Sai contract with VulcanizeDB

## Setup
Setup Postgres and Geth - see [VulcanizeDB README](https://github.com/vulcanize/VulcanizeDB/blob/master/README.md)

```
make installtools
make setup NAME=vulcanize_public
make migrate NAME=vulcanize_public
make build
```

`./sai_watcher sync --config environments/public.toml --starting-block-number n` (where n is a recent enough block to fetch contract data with your node)

`./sai_watcher getEvents --config environments/public.toml`

## Graphql
Compatible with `postgraphile`:

```
npm install -g postgraphile
postgraphile -c "postgresql://user@localhost:5432/vulcanize_public" --schema=public,maker
```

Example schema / server in `go` in `graphql_server/schema.go`:

```
./sai_watcher graphql --config environments/public.toml
```

## IPFS Reports
Create a report of your database and persist it on IPFS:

Install IPFS

```
./sai_watcher generateReport --config environments/public.toml
```

## Run tests
```
make installtools
make setup NAME=vulcanize_private
make migrate NAME=vulcanize_private
ginkgo -r
```
