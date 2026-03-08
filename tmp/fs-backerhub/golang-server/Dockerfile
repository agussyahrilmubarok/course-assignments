FROM alpine:3.22
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY bin/backerhub-golang-server /app/backerhub-golang-server
COPY configs/config-stage.json /app/configs/config.json
COPY public /app/public

EXPOSE 8080
# ENTRYPOINT ["/app/backerhub-golang-server", "-config", "/app/config.prod.json"]
ENTRYPOINT ["/app/backerhub-golang-server"]