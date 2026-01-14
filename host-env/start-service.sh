#!/usr/bin/env bash

BASE_DIR="$HOME/PANEL"
APP_DIR="$BASE_DIR/app"
mkdir -p "$BASE_DIR" "$APP_DIR"
cp ./* "$APP_DIR/"
cp ../filter-power/filter-power "$APP_DIR/"

cd "$APP_DIR" || exit

sudo chmod +x get-device-info.sh upload-power.sh get-device-upload.sh filter-power

# Systemd services
OPTION=$1

if [ $OPTION -eq "reload" ]; then
    sudo cp device-upload.service device-upload.timer /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable device-upload.timer
    sudo systemctl restart device-upload.timer
fi
