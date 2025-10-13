#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit

export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_UPLOAD_PATH="/home/panel/csv-save"
export SERVER_IP="173.212.201.115"
export SERVER_PORT="5040"
export IMEI_MAP="/home/panel/id2imei.csv"

mkdir -p CSV_INPUT_PATH CSV_UPLOAD_PATH

./get-device-info.sh || {
  echo "Download device information failed"; exit 1
}
./filter-power $CSV_INPUT_PATH $CSV_UPLOAD_PATH || {
  echo "Filter and upload power data failed"; exit 1
}

