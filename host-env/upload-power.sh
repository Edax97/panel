#!/bin/bash

export SERVER_IP
export SERVER_PORT
export IMEI_MAP
export TRANSFORMER_IDS
SERVER_IP="$(bws secret get 7c7626f4-e27b-4eac-92c5-b37500fd60c6 | jq -r '.value')"
SERVER_PORT="$(bws secret get b48289e0-aef1-4146-bb2e-b37500fd74aa | jq -r '.value')"
IMEI_MAP="$(bws secret get 3067744e-5edd-4568-8b45-b37500ff07b0 | jq -r '.value')"
TRANSFORMER_IDS="$(bws secret get b35a11ae-f97b-44ea-acc8-b39f016ea87b | jq -r '.value')"

CSV_INPUT="$1"
CSV_FILTER="$2"
if [ -z "$CSV_INPUT" ] || [ -z "$CSV_FILTER" ]; then
  echo "CSV_INPUT|UPLOAD_PATH not set" >&2
  exit 1
fi
mkdir -p "$CSV_FILTER"

./filter-power "$CSV_INPUT" "$CSV_FILTER" || {
  echo "Filter and upload power data failed"; exit 1
}
