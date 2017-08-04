#!/bin/sh

echo "local bin"

ls -la ~/.local/bin/

# Check that the environment variable has been set correctly
if [ -z "$SECRETS_BUCKET_NAME" ]; then
  echo >&2 'error: missing SECRETS_BUCKET_NAME environment variable'
  exit 1
fi

mkdir -p config
# Load the S3 secrets

~/.local/bin/aws s3 cp s3://"${SECRETS_BUCKET_NAME}"/config.yaml config/

./lvl-api
