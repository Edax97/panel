#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit
# curl github/env-file
# source env-file
export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_SAVE_PATH="/home/panel/csv-save"

./bin-get
./bin-save
/usr/bin/podman pod start cloud-pod

