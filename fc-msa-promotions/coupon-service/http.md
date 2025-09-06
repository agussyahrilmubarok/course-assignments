# Create CouponPolicy

`POST /api/couponPolicies`

```bash
curl -X POST http://localhost:8082/api/couponPolicies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Coupon",
    "description": "Get 10% off for all items",
    "discountType": "PERCENTAGE",
    "discountValue": 10,
    "minimumOrderAmount": 50000,
    "maximumDiscountAmount": 20000,
    "totalQuantity": 100,
    "startTime": "2025-01-01T00:00:00",
    "endTime": "2025-12-31T23:59:59"
  }' \
  -i
```

# Get Coupon Policies

`GET /api/couponPolicies`

```bash
curl -X GET http://localhost:8082/api/couponPolicies \
  -H "Content-Type: application/json" \
  -i
```

# Get Coupon Policy

`GET /api/couponPolicies/{id}`

```bash
curl -X GET http://localhost:8082/api/couponPolicies/ID \
  -H "Content-Type: application/json" \
  -i
```

# Issue Coupon

`POST /api/coupons/issue`

```bash
curl -X POST http://localhost:8082/api/coupons/issue \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "couponPolicyId": "ID"
  }' \
  -i
```

# Use Coupon

`POST /api/coupons/{ID}/use`

```bash
curl -X POST http://localhost:8082/api/coupons/ID/use \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -d '{
    "orderId": "ORDER-1X"
  }' \
  -i
```

# Cancel Coupon

`POST /api/coupons/{ID}/cancel`

```bash
curl -X POST http://localhost:8082/api/coupons/ID/cancel \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```

# Get Coupons

`GET /api/coupons`

```bash
curl -X GET http://localhost:8082/api/coupons \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```

# Get Coupon

`GET /api/coupons/{ID}`

```bash
curl -X GET http://localhost:8082/api/coupons/ID \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: 10001" \
  -i
```