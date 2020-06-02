FROM alpine:latest

RUN mkdir -p /etc/multicloud
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY  ./pkg/examples/policy.json /etc/multicloud/
COPY api /api
ENTRYPOINT ["/api"]
