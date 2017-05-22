#!/usr/bin/env bash
# run ramped RPS test against an endpoint
set -e

target=${1:-hello}
file="results/${target}.csv"

rm $file || true
for rps in $(seq 400 5 600); do
    ./x-files-rate -target $target -rps $rps | tee -a $file
done
echo 'e' >> $file

