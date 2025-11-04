#!/bin/bash
FILE_TG=/home/edmar/Documents/mall-trujillo/history-all/t1_2/SERVERTG1-TG2.csv
IMEI_FILE=/home/edmar/Documents/panel/map-imei/map-imei.csv
IDS_FILE="$1"

row="$( awk 'NR==4' "$FILE_TG" )"
echo "ROW IMEI> $row"
IFS=";"
read -ra imeiArray <<< "$row"
unset IFS

ID_CSV=""
re='^[0-9]+$'
for imei in "${imeiArray[@]}"; do
  echo "Item: $imei"
  if ! [[ $imei =~ $re ]]; then
    ID="."
  else
    ID="$( grep "$imei" "$IMEI_FILE" | awk -F',' '{print $1}' | grep -oP '((?<=csv_)).*' )"
  fi
  ID_CSV="$ID_CSV;$ID"
#  echo "$ID_CSV"
done
rm "$IDS_FILE"
printf "%s\n" "$ID_CSV" > "$IDS_FILE"

