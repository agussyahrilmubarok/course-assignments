FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY bin/gocoupon-api /app/api
COPY config-prod.yml /app/config.yml
RUN chmod +x /app/api
RUN chmod +x /app/config.yml

EXPOSE 8080

ENTRYPOINT ["/app/api"]