#!/usr/bin/env bash
mkdir -p /home/panel/app
cp ./* /home/panel/app
cd /home/panel/app || exit
go build -o filter-power /home/panel/panel/filter-power
sudo chmod +x get-device-info.sh upload-power.sh
mkdir -p csv-input csv-save

# Systemd services
cp save-energy.service save-energy.timer .config/systemd/user/
systemctl --user daemon-reload
systemctl enable --user save-energy.timer
systemctl start --user save-energy.timer