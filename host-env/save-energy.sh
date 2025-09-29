#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit

export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_SAVE_PATH="/home/panel/csv-save"
mkdir -p CSV_INPUT_PATH CSV_SAVE_PATH

./get-panels.sh || {
  echo "Get panel data failed"; exit 1
}
./bin-save || {
  echo "Save csv failed"; exit 1
}
/usr/bin/podman pod start cloud-pod

