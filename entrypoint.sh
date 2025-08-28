#!/bin/sh
set -e

echo "[DEBUG] Current environment:"
env | grep -E 'ENVIRONMENT|GAMERESULT_PORT|GAMERESULT_URL|GATESERVER_PORT'

# Render the config file
envsubst < /app/config.template.yaml > /app/config.yaml

echo "[INFO] Rendered config.yaml:"
cat /app/config.yaml

# Run your Go app
exec ./main