# Sai Contract Watcher

## Description
A [VulcanizeDB](https://github.com/vulcanize/VulcanizeDB) transformer for watching events related to the Sai contract.

## Dependencies
 - Go 1.9+
 - Postgres 10
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)
 - [IPFS](https://github.com/ipfs/go-ipfs#build-from-source)

## Installation
1. Setup Postgres and an Ethereum node - see [VulcanizeDB README](https://github.com/vulcanize/VulcanizeDB/blob/master/README.md).
1. `git clone git@github.com:8thlight/sai_watcher.git`

  _note: `go get` does not work for this project because need to run the (fixlibcrypto)[https://github.com/8thlight/sai_watcher/blob/master/Makefile] command along with `go build`._
1. Install dependencies:
    ```
    make installtools
    ```
1. Create the database based on the [VulcanizeDB schema](https://github.com/vulcanize/VulcanizeDB/blob/master/db/schema.sql):
    ```
    make setup NAME=vulcanize_public
    ```
1. Run the migrations to add project specific tables to the database:
    ```
    make migrate NAME=vulcanize_public
    ```
1. Build:
    ```
    make build
    ```

## Configuration
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `ipcPath` to the local node's IPC filepath:
  - when using geth:
    - The IPC file is called `geth.ipc`.
    - The geth IPC file path is printed to the console when you start geth.
    - The default location is:
      - Mac: `$HOME/Library/Ethereum`
      - Linux: `$HOME/.ethereum`

  - when using parity:
    - The IPC file is called `jsonrpc.ipc`.
    - The default location is:
      - Mac: `$HOME/Library/Application\ Support/io.parity.ethereum/`
      - Linux: `$HOME/.local/share/io.parity.ethereum/`

- See `environments/infura.toml` to configure commands to run against infura, if a local node is unavailable.

## Running the sync command
This command syncs VulcanizeDB with the configured Ethereum node.
1. Start node (**if fast syncing wait for initial sync to finish**)
1. In a separate terminal window:
  `./sai_watcher sync --config <config.toml> --starting-block-number <block-number>`
  - where `block-number` is a recent enough block to fetch contract data with your node

## Running the getEvents command
`getEvents` starts up a process to watch for blocks on specified contracts, as well as specific log events associated with those contracts. It then stores transformed values in the following tables and views in the VulcanizeDB database:

Tables:
- `maker.peps_everyblock`
- `maker.cup_action`
- `log_filters`
- `maker.gov`

Views:
- `public.cup_act`
- `public.cup`

This command will need to be run against a full archive node. If a local full archive node is unavailable, see the previous point about running
this command against infura.

`./sai_watcher getEvents --config <config.toml>`

## Graphql
We're using [PostGraphile](https://www.graphile.org/postgraphile/) to create a GraphQL API from the VulcanizeDB postgres schema.
1. Ensure that Node.js v8.6 is installed.
1. Install postgraphile:
    ```
    npm install -g postgraphile
    ```
1. Start the postgraphile server:
    ```
    postgraphile -c "postgresql://<user>@localhost:5432/vulcanize_public" --schema=public,maker
    ```
    - the `-c "postgresql://user@localhost:5432/vulcanize_public"` flag indicates which postgres connection postgraphile should be looking to, where `<user>` is your local postgres user
    - the `--schema=public,maker` flag indicates which schema(s) postgraphile should use to generate the GraphQL API

## IPFS Reports
This task creates a report of your database and persists it on IPFS.
1. Run the command:
    ```
    ./sai_watcher generateReport --config environments/public.toml
    ```
1. This will return an IPFS hash.
1. To fetch this hash from IPFS:
    ```
    ipfs dag get <the hash>
    ```

## Running the tests
```
make installtools
make setup NAME=vulcanize_private
make migrate NAME=vulcanize_private
ginkgo -r
```
