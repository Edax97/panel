#!/usr/bin/env bash

login (){
  local URL="$1/em-edm/sessions"
  local USER="$2"
  local PASS="$3"
  local PAYLOAD="{\"scheme\":\"BASIC\",\"user\":\"$USER\",\"password\":\"$PASS\"}"
  local res="$(curl -k -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD")"
  jq -r '.AccessToken."access_token"' <<< "$res"
}
download (){
  local URL="$1"
  local TOKEN="$2"
  local FILENAME="$3"
  local FROM
  local TO
  FROM="$(date --date='1 hour ago' +%Y-%m-%dT%H:%M)"
  TO="$(date +%Y-%m-%dT%H:%M)"
  curl -k "$URL/csv?start_date=$FROM&end_date=$TO" \
    -H "Accept: text/csv, */*" \
    -H "Accept-Encoding: gzip, deflate" \
    -H "Connection: keep-alive" \
    -H "Authorization:  Bearer $TOKEN" \
    -H "Referer: $URL/public/settings/equipment-management/local-export" \
    -H 'Sec-Fetch-Dest: empty' \
    -H 'Sec-Fetch-Mode: cors' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Sec-GPC: 1' \
    --compressed \
    --output "$FILENAME" --fail
}

PANEL_USERS="$(bws secret get f2db263d-b244-482d-a37f-b3640162669d | jq -r '.value')"
PANEL_PASS="$(bws secret get 62b15eba-9580-4792-8658-b36401628dc4 | jq -r '.value')"
PANEL_URLS="$(bws secret get 4aa7f245-6764-4820-8ebb-b3640161a5c7 | jq -r '.value')"

readarray -t urls <<< "$PANEL_URLS"
readarray -t psws <<< "$PANEL_PASS"
readarray -t users <<< "$PANEL_USERS"

if [ -z "$CSV_INPUT_PATH" ]; then
  echo "CSV_INPUT_PATH not set" >&2
  exit 1
fi

rm -r "$CSV_INPUT_PATH" 2>/dev/null || echo "Creating dir..."
mkdir -p "$CSV_INPUT_PATH"

i=0
if [ "${#urls[@]}" -eq "${#psws[@]}" ] && [ "${#urls[@]}" -eq "${#users[@]}" ] && [ "${#urls[@]}" -gt 0 ]; then
  for i in "${!urls[@]}"; do
    url="${urls[i]}"
    pass="${psws[i]}"
    user="${users[i]}"
    echo ">> Log in to $url"
    if [ -n "$url" ] && [ -n "$user" ] && [ -n "$pass" ]; then
      token=$(login "$url" "$user" "$pass")
      if [ -n "$token" ]; then
        echo ">> Downloading device data $(( i+1 ))..."
        download "$url" "$token" "$CSV_INPUT_PATH/data_$(( i+1 )).csv"
      else
        echo "No se pudo obtener token: $url" >&2
      fi
    fi
  done
else
  echo "Secretos tienen distinto numero de lineas" >&2
  exit 1
fi


