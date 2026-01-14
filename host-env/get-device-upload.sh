#!/bin/bash

source /home/pi/.env
cd "$HOME/PANEL/app" || exit

export CSV_INPUT_PATH="$HOME/PANEL/store/input"
export CSV_UPLOAD_PATH="$HOME/PANEL/store/save"

mkdir -p "$CSV_INPUT_PATH" "$CSV_UPLOAD_PATH"

./get-device-info.sh "$CSV_INPUT_PATH" || {
  echo "Download device information failed"
}
./upload-power.sh "$CSV_INPUT_PATH" "$CSV_UPLOAD_PATH" || {
  echo "Filter and upload power data failed"; exit 1
}
