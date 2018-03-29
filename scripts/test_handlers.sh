#! /bin/bash
set -e
echo "Setting up Vulcanize"
trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT
make setup NAME=vulcanize_public
make migrate NAME=vulcanize_public
make build

./sai_watcher sync --config environments/public.toml --starting-block-number 4752014 & 
echo "Running Watcher"
./sai_watcher getEvents --config environments/infura.toml 
