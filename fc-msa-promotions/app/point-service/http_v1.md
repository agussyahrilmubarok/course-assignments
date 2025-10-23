# Earn Point

`POST /api/v1/points/earn`

```bash
curl -X POST http://localhost:8083/api/v1/points/earn \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "amount": 10,
    "description": "Get 10 point for some promo"
  }' \
  -i
```

# Use Point

`POST /api/v1/points/use`

```bash
curl -X POST http://localhost:8083/api/v1/points/use \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "amount": 5,
    "description": "Use 5 point for some item"
  }' \
  -i
```

# Cancel Point

`POST /api/v1/points/cancel`

```bash
curl -X POST http://localhost:8083/api/v1/points/cancel \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "pointId": "",
    "description": "Cancel point"
  }' \
  -i
```

# Get Balance Point

`POST /api/v1/points/users/balance`

```bash
curl -X GET http://localhost:8083/api/v1/points/users/balance \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```

# Get History Point

`POST /api/v1/points/users/history`

```bash
curl -X GET http://localhost:8083/api/v1/points/users/history \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```