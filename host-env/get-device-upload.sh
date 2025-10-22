#!/usr/bin/env bash
cd /home/panel/app || exit

export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_UPLOAD_PATH="/home/panel/csv-save"

mkdir -p "$CSV_INPUT_PATH" "$CSV_UPLOAD_PATH"

./get-device-info.sh "$CSV_INPUT_PATH" || {
  echo "Download device information failed"
}
./upload-power.sh "$CSV_INPUT_PATH" "$CSV_UPLOAD_PATH" || {
  echo "Filter and upload power data failed"; exit 1
}