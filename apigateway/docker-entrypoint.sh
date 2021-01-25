#!/bin/sh
set -e

echo "Generating envoy.yaml config file..."
cat /tmp/envoy.yaml.tmpl | envsubst \$BOOKING_SVC_PORT\$TRACKING_SVC_PORT\$HANDLING_SVC_PORT > /etc/envoy.yaml

echo "Starting Envoy..."
/usr/local/bin/envoy -c /etc/envoy.yaml --bootstrap-version 2