# Create Time Sale

`/api/v3/timeSale`

```bash
curl -X POST http://localhost:8084/api/v3/timeSales \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "product": {
      "name": "Gaming Mouse",
      "price": 499000
    },
    "quantity": 100,
    "discountPrice": 299000,
    "startAt": "2025-10-11T00:00:00",
    "endAt": "2025-12-31T23:59:59"
  }' \
  -i
```

# Get Time Sale

`/api/v3/timeSale/{id}`

```bash
curl -X GET http://localhost:8084/api/v3/timeSales/ \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```

# Get Ongoing Time Sale

`/api/v3/timeSale`

```bash
curl -X GET http://localhost:8084/api/v3/timeSales \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```

# Purchase Time Sale

`/api/v3/timeSale/purchase`

```bash
curl -X POST http://localhost:8084/api/v3/timeSales/purchase \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "timeSaleId": "",
    "quantity": 50
  }' \
  -i
```

# Get Purchase Time Sale

`/api/v3/timeSale/{timeSaleId}/purchase/results/{requestId}`

```bash
curl -X GET http://localhost:8084/api/v3/timeSales/{timeSaleId}/purchase/results/{requestId} \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```