# Earn Point

`POST /api/v2/points/earn`

```bash
curl -X POST http://localhost:8083/api/v2/points/earn \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10002" \
  -d '{
    "amount": 10,
    "description": "Get 10 point for some promo"
  }' \
  -i
```

# Use Point

`POST /api/v2/points/use`

```bash
curl -X POST http://localhost:8083/api/v2/points/use \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10002" \
  -d '{
    "amount": 5,
    "description": "Use 5 point for some item"
  }' \
  -i
```

# Cancel Point

`POST /api/v2/points/cancel`

```bash
curl -X POST http://localhost:8083/api/v2/points/cancel \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10002" \
  -d '{
    "pointId": "54efe237-b646-4ad6-b831-d85b3cf08c4a",
    "description": "Cancel point"
  }' \
  -i
```

# Get Balance Point

`POST /api/v2/points/users/balance`

```bash
curl -X GET http://localhost:8083/api/v2/points/users/balance \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10002" \
  -i
```

# Get History Point

`POST /api/v2/points/users/history`

```bash
curl -X GET http://localhost:8083/api/v2/points/users/history \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10002" \
  -i
```