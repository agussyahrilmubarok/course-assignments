FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY bin/api /app/api
COPY config-prod.yml /app/config.yml

EXPOSE 8080

ENTRYPOINT ["/app/api"]