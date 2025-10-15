#!/usr/bin/env bash

DOWN_DIR="$1"
if [ -z "$DONW_DIR" ]; then
  DOWN_DIR="."
fi
aws s3 cp "s3://cloud-store-mall/data-export_20251014 (tg2n1).csv" "$DOWN_DIR"
aws s3 cp "s3://cloud-store-mall/data-export_20251014 (tg2n2).csv" "$DOWN_DIR"
aws s3 cp "s3://cloud-store-mall/data-export_20251015 (electrogeno).csv" "$DOWN_DIR"