#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit

export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_SAVE_PATH="/home/panel/csv-save"
export SERVER_IP="173.212.201.115"
export SERVER_PORT="5040"
export IMEI_MAP="/home/panel/id2imei.csv"

mkdir -p CSV_INPUT_PATH CSV_SAVE_PATH

./get-panels.sh || {
  echo "Download panel data failed"; exit 1
}
./filter-csv $CSV_INPUT_PATH $CSV_SAVE_PATH || {
  echo "Filter csv failed"; exit 1
}

