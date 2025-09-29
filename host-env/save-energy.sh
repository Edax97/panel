#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit


export PANEL_URLS PANEL_PASS PANEL_USERS
export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_SAVE_PATH="/home/panel/csv-save"

./bin-get || {
  echo "Get panel data failed"; exit 1
}
./bin-save || {
  echo "Save csv failed"; exit 1
}
/usr/bin/podman pod start cloud-pod

