# Register Product

`api/v1/catalogs/product`

```bash
curl -X POST http://localhost:8082/api/v1/catalogs/products \
  -H "Content-Type: application/json" \
  -d '{
    "sellerId": "seller-abc-123",
    "name": "Wireless Bluetooth Headphones",
    "description": "Noise-cancelling over-ear headphones with 40 hours of battery life.",
    "price": 249000,
    "stockCount": 50,
    "tags": ["electronics", "audio", "wireless", "headphones"]
  }' \
  -i
```

# Decrease Stock Count

`api/v1/catalogs/products/{productId}/decreaseStockCount`

```bash
curl -X POST http://localhost:8082/api/v1/catalogs/products//decreaseStockCount \
  -H "Content-Type: application/json" \
  -d '{
    "decreaseCount": 2
  }' \
  -i
```

# Get Product By Id

`api/v1/catalogs/products/{productId}`

```bash
curl -X GET http://localhost:8082/api/v1/catalogs/products/ \
  -H "Content-Type: application/json" \
  -i
```

# Get Products By Seller Id

`api/v1/catalogs/sellers/{sellerId}/products`

```bash
curl -X GET http://localhost:8082/api/v1/catalogs/sellers/seller-abc-123/products \
  -H "Content-Type: application/json" \
  -i
```

# Delete Product By Id

`api/v1/catalogs/products/`

```bash
curl -X DELETE http://localhost:8082/api/v1/catalogs/products/1ca61ddb-0db9-41cd-ac1f-86d1f799f013 \
  -H "Content-Type: application/json" \
  -i
```