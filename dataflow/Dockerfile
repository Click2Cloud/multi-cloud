FROM alpine:latest

COPY dataflow /dataflow
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/dataflow"]
