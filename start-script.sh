#!/usr/bin/env bash

rm "$STORAGE_DIR/uploaded" || "Already deleted /uploaded"
cd /app || exit
mkdir -p "${STORAGE_DIR || "-"}"
./bin "$STORAGE_DIR" "$URL"
touch "$STORAGE_DIR/uploaded"
