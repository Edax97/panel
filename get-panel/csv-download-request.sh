#!/usr/bin/env bash
URL="https://192.168.8.104"
PATH="csv?start_date=2025-09-01T00:00&end_date=2025-09-01T02:00"
TOKEN="$(jq '.AccessToken.access_token' login-response.json)"

curl -k "$URL/$PATH" \
    -H "Accept: text/csv, /" \
    -H "Authorization:  Bearer $TOKEN" \
    -H 'Connection: keep-alive' \
    -H "Referer: $URL/public/settings/equipment-management/local-export" \
    -H 'Sec-Fetch-Dest: empty' \
    -H 'Sec-Fetch-Mode: cors' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Sec-GPC: 1'