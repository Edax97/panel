#!/bin/bash
FILE_TO_UPLOAD="/storage"
sleep 5

echo "Starting upload... $BUCKET"
export AWS_ACCESS_KEY_ID=$ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=$SECRET_KEY

/usr/local/bin/aws s3 cp "$FILE_TO_UPLOAD" "s3://$BUCKET" --recursive --region "$REGION"
