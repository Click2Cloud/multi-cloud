FROM alpine:latest

COPY datamover /datamover
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/datamover"]
