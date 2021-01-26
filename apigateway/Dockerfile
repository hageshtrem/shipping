FROM envoyproxy/envoy:v1.17.0

RUN apt-get update && \
    apt-get install gettext -y

WORKDIR /apigateway

COPY envoy.yaml envoy.yaml
COPY descriptor.pb descriptor.pb
COPY docker-entrypoint.sh docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]