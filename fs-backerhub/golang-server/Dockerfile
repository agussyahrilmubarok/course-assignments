FROM golang:1.25.3 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backerhub cmd/server/main.go

FROM alpine:latest
# ENV GIN_MODE=release
WORKDIR /app
COPY --from=builder /build /app
COPY configs/backerhub-stage.json /app/configs/backerhub.json
COPY public /app/public

EXPOSE 8080

# ENTRYPOINT ["/app/backerhub", "-config", "/app/config.prod.json"]
ENTRYPOINT ["/app/backerhub"]