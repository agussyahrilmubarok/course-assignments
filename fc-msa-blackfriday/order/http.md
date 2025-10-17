# Start Order

`api/v1/orders/start`

```bash
curl -X POST http://localhost:8084/api/v1/orders/start \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-1-x",
    "productId": "24cb1090-3f46-438b-898e-743c8e056bcf",
    "count": 5
  }' \
  -i
```