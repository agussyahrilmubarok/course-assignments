FROM alpine:3.22
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY bin/hello-docker /app/hello-docker
ENTRYPOINT ["/app/hello-docker"]