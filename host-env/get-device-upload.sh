#!/usr/bin/env bash
cd ~/app || exit

export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_UPLOAD_PATH="/home/panel/csv-save"

mkdir -p "$CSV_INPUT_PATH" "$CSV_UPLOAD_PATH"

./get-device-info.sh || {
  echo "Download device information failed"; exit 1
}
./upload-power.sh || {
  echo "Filter and upload power data failed"; exit 1
}