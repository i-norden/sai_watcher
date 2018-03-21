#! /bin/bash
set -e
echo "Setting up Vulcanize"
VULCANIZE_DIR=~/go/src/github.com/vulcanize/vulcanizedb
trap 'kill $(jobs -pr)' SIGINT SIGTERM EXIT
${VULCANIZE_DIR}/test_scripts/fresh_vulcanize.sh
psql vulcanize_public -c "drop table schema_migrations;"
make migrate NAME=vulcanize_public

go build
./sai_watcher sync --config environments/public.toml --starting-block-number 4752014 &
sleep 15
./sai_watcher getEvents --config environments/infura.toml 
