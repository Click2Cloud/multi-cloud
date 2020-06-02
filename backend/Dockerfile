FROM alpine:latest

COPY backend /backend
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/backend"]
