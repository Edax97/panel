#!/usr/bin/env bash

login (){
  local URL="$1/em-edm/sessions"
  local USER="$2"
  local PASS="$3"
  local PAYLOAD="{\"scheme\":\"BASIC\",\"user\":\"$USER\",\"password\":\"$PASS\"}"
  curl -k -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD" | jq '.AccessToken."access-token"'
}
download (){
  local URL="$1"
  local TOKEN="$2"
  local DEST_PATH="$3"
  echo "Saving to $DEST_PATH"
  curl -k "$URL/csv?start_date=2025-09-01T00:00&end_date=2025-09-01T02:00" \
    -H "Accept: text/csv, /" \
    -H "Authorization:  Bearer $TOKEN" \
    -H 'Connection: keep-alive' \
    -H "Referer: $URL/public/settings/equipment-management/local-export" \
    -H 'Sec-Fetch-Dest: empty' \
    -H 'Sec-Fetch-Mode: cors' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Sec-GPC: 1' > "$DEST_PATH"
}

PANEL_USERS="$(bws secret get f2db263d-b244-482d-a37f-b3640162669d | jq -r '.value')"
PANEL_PASS="$(bws secret get 62b15eba-9580-4792-8658-b36401628dc4 | jq -r '.value')"
PANEL_URLS="$(bws secret get 4aa7f245-6764-4820-8ebb-b3640161a5c7 | jq -r '.value')"
export CSV_INPUT_PATH="/home/panel/csv-input"
export CSV_SAVE_PATH="/home/panel/csv-save"

readarray -t urls <<< "$PANEL_URLS"
readarray -t psws <<< "$PANEL_PASS"
readarray -t users <<< "$PANEL_USERS"

i=0
if [ "${#urls[@]}" -eq "${#psws[@]}" ] && [ "${#urls[@]}" -eq "${#users[@]}" ] && [ "${#urls[@]}" -gt 0 ]; then
  for i in "${!urls[@]}"; do
    url="${urls[i]}"
    pass="${psws[i]}"
    user="${users[i]}"
    echo "- $url $user $pass"
    if [ -n "$url" ] && [ -n "$user" ] && [ -n "$pass" ]; then
      token=$(login "$url" "$user" "$pass")
      if [ -n "$token" ]; then
        download "$url" "$token" "$CSV_INPUT_PATH/data_$i.csv"
      else
        echo "No se pudo obtener token: $url" >&2
      fi
    fi
  done
else
  echo "Secretos tienen distinto numero de lineas" >&2
  exit 1
fi


