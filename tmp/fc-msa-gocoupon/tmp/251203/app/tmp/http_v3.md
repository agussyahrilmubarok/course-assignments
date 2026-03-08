# HTTP v3 Example

## Init Dummy Coupon Policy

```bash
curl -X GET http://localhost:8080/init-dummy-redis-db \
  -H "Content-Type: application/json" \
  -i
```

## Clean Dummy Coupon Policy

```bash
curl -X GET http://localhost:8080/clean-dummy-redis-db \
  -H "Content-Type: application/json" \
  -i
```

## Issue Coupon Request V2

```bash
curl -X POST http://localhost:8080/api/v3/coupons/issue \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "policy_code": "BF-C10"
  }' \
  -i
```

## Issue Coupon Request V2 Loop (11 requests)

```bash
for i in {1..11}
do
  USER_ID="USER_$i"
  echo "Request ke-$i dengan X-USER-ID: $USER_ID"

  curl -X POST http://localhost:8080/api/v3/coupons/issue \
    -H "Content-Type: application/json" \
    -H "X-USER-ID: $USER_ID" \
    -d '{"policy_code": "BF-C10"}' \
    -i

  echo "" 
done
```

## Use Coupon Request V2

```bash
curl -X POST http://localhost:8080/api/v3/coupons/use \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "coupon_code": "",
    "order_id": "ORDER-12345"
  }' \
  -i
```

## Cancel Coupon Request V2

```bash
curl -X POST http://localhost:8080/api/v3/coupons/cancel \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "coupon_code": ""
  }' \
  -i
```

## Find Coupon By Code

```bash
curl -X GET http://localhost:8080/api/v3/coupons/dab671d9-8b39-453c-bf13-0d74b66c8e1e \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -i
```
