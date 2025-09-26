#!/usr/bin/env bash
cd /home/panel || exit
mkdir -p csv-input csv-save
podman kube play --replace pod.yaml
cp save-energy.service save-energy.timer .config/systemd/user/
systemctl --user daemon-reload
systemctl enable --user save-energy.timer
systemctl start --user save-energy.timer