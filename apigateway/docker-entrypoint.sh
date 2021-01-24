#!/bin/sh
set -e

echo "Generating envoy.yaml config file..."
cat /tmp/envoy.yaml.tmpl | envsubst \$TRACKING_SVC_PORT > /etc/envoy.yaml

echo "Starting Envoy..."
/usr/local/bin/envoy -c /etc/envoy.yaml --bootstrap-version 2