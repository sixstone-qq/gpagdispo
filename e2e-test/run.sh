#!/bin/bash

set -e

source "e2e-test/assert.sh"

trap stopSvcs EXIT

pushd checker
go build ./cmd/gpagdispo-checker/...
CONFIG_PATH=websites.ion ./gpagdispo-checker &
checkerPID=$!
popd

# Require to get into directory to find migration files
pushd recorder
go build ./cmd/gpagdispo-recorder/...
./gpagdispo-recorder &
recorderPID=$!
popd

# Wait for 10s to produce data
echo "Sleep for 10s..."
sleep 10

stopSvcs() {
    kill $checkerPID $recorderPID || true
}

stopSvcs

websites=$(docker exec gpagdispo_postgres_1 psql -U postgres -h localhost website_monitor -t -c 'SELECT COUNT(*) FROM websites;' | xargs)
results=$(docker exec gpagdispo_postgres_1 psql -U postgres -h localhost website_monitor -t -c 'SELECT COUNT(*) FROM websites_results' | xargs)

assert "$websites" -eq 3
assert "$results" -gt 6
