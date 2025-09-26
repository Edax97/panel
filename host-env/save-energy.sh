#!/usr/bin/env bash
mkdir -p ~/app
cd ~/app || exit


PANEL_USERS="$(bws secret get f2db263d-b244-482d-a37f-b3640162669d | jq '.value')"
PANEL_PASS="$(bws secret get 62b15eba-9580-4792-8658-b36401628dc4 | jq '.value')"
PANEL_URLS="$(bws secret get 4aa7f245-6764-4820-8ebb-b3640161a5c7 | jq '.value')"
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

