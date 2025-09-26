#!/usr/bin/env bash
cp ./* /home/panel/
cd /home/panel || exit
go build -o bin-save ./panel/save-csv
go build -o bin-get ./panel/get-panel
mkdir -p csv-input csv-save
podman kube play --replace pod.yaml
cp save-energy.service save-energy.timer .config/systemd/user/
systemctl --user daemon-reload
systemctl enable --user save-energy.timer
systemctl start --user save-energy.timer