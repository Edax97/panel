#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit

export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_UPLOAD_PATH="/home/panel/csv-save"
export SERVER_IP
export SERVER_PORT
export IMEI_MAP
SERVER_IP="$(bws secret get 7c7626f4-e27b-4eac-92c5-b37500fd60c6 | jq -r '.value')"
SERVER_PORT="$(bws secret get b48289e0-aef1-4146-bb2e-b37500fd74aa | jq -r '.value')"
IMEI_MAP="$(bws secret get 3067744e-5edd-4568-8b45-b37500ff07b0 | jq -r '.value')"
# export SERVER_PORT="5040"
# export IMEI_MAP="/home/panel/id2imei.csv"

rm -r "$CSV_INPUT_PATH"
mkdir -p "$CSV_INPUT_PATH" "$CSV_UPLOAD_PATH"

./get-device-info.sh || {
  echo "Download device information failed"; exit 1
}
./filter-power $CSV_INPUT_PATH $CSV_UPLOAD_PATH || {
  echo "Filter and upload power data failed"; exit 1
}

