FROM envoyproxy/envoy:v1.17.0

RUN apt-get update && \
    apt-get install gettext -y

COPY envoy.yaml /tmp/envoy.yaml.tmpl
COPY descriptor.pb /tmp/envoy/descriptor.pb
COPY docker-entrypoint.sh /

RUN chmod 500 /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]