#!/usr/bin/env bash
cd /home/panel || exit
mkdir -p storage
podman kube play --replace pod.yaml
podman generate systemd -f --name panel-pod.yaml
cp *.service panel-pod.timer .config/systemd/user/
systemctl --user daemon-reload
systemctl enable --user panel-pod.timer
systemctl start --user panel-pod.timer