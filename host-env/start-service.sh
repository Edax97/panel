#!/usr/bin/env bash
APP_DIR="/home/panel/app"
mkdir -p "$APP_DIR"
cp ./* "$APP_DIR"
cd "$APP_DIR" || exit

cd /home/panel/panel/filter-power || exit
go build -o "$APP_DIR/filter-power" .
cd "$APP_DIR" || exit

sudo chmod +x get-device-info.sh upload-power.sh
mkdir -p csv-input csv-save

# Systemd services
#cp device-upload.service device-upload.timer .config/systemd/user/
#systemctl --user daemon-reload
#systemctl enable --user device-upload.timer
#systemctl start --user device-upload.timer