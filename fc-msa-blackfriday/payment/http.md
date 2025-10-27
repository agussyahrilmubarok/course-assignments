# Create Transaction

`api/v1/transactions`

```bash
curl -X POST http://localhost:8085/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "orderId-123",
    "amount": 249000
  }' \
  -i
```

# Get Transaction

`api/v1/transactions/{id}`

```bash
curl -X GET http://localhost:8085/api/v1/transactions/sss \
  -H "Content-Type: application/json" \
  -i
```