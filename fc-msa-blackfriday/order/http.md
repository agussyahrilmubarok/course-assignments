# Start Order

`api/v1/orders/create`

```bash
curl -X POST http://localhost:8084/api/v1/orders/create \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-1-x",
    "productId": "08d7de89-9f8a-42d5-866c-c4440116e30a",
    "count": 5
  }' \
  -i
```

# Find Order

`api/v1/orders/{id}`

```bash
curl -X GET http://localhost:8084/api/v1/orders/a9590ad9-4bb6-4e33-815b-7e0db2f48673 \
  -H "Content-Type: application/json" \
  -i
```

# Find Order

`api/v1/orders/users/{id}`

```bash
curl -X GET http://localhost:8084/api/v1/orders/users/user-1-x \
  -H "Content-Type: application/json" \
  -i
```