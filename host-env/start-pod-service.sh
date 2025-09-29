#!/usr/bin/env bash
cp ./* /home/panel/
cd /home/panel || exit
go build -o bin-save ./panel/save-csv
sudo chmod +x ./get-panels.sh
mkdir -p csv-input csv-save
podman kube play --replace pod.yaml
# Systemd services
cp save-energy.service save-energy.timer .config/systemd/user/
systemctl --user daemon-reload
systemctl enable --user save-energy.timer
systemctl start --user save-energy.timer