#!/usr/bin/env bash
URL="https://192.168.8.104/em-edm/sessions"
PAYLOAD='{"scheme":"BASIC","user":"SecurityAdmin","password":"UGFuZWxvY2F0YXJpbzEu"}'
curl -kv -X POST "$URL" \
-H "Content-Type: application/json" \
-d "$PAYLOAD"