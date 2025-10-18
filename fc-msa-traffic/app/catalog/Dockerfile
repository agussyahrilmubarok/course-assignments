FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o catalog-service ./cmd/server/main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/catalog-service .
COPY config-prod.yaml /app/configs/config.yaml

EXPOSE 8082
ENTRYPOINT ["./catalog-service"]
