# Graylog Server

## Run

1. Compose Run

```bash
docker compose up -d --build
```

2. Setup inputs on http://localhost:9000

```bash
# System/Inputs > Inputs > GELF UDP > Launch new input
Title: GELFUDP
BindAddress: 0.0.0.0
Port: 12201

```