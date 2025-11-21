# HTTP Example

## Init Dummy Coupon Policy V1

```bash
curl -X GET http://localhost:8080/init-dummy-v1 \
  -H "Content-Type: application/json" \
  -i
```

## Init Dummy Coupon Policy V2

```bash
curl -X GET http://localhost:8080/init-dummy-v2 \
  -H "Content-Type: application/json" \
  -i
```

## Clean Dummy Coupon Policy V1

```bash
curl -X GET http://localhost:8080/clean-dummy-v1 \
  -H "Content-Type: application/json" \
  -i
```

## Clean Dummy Coupon Policy V2

```bash
curl -X GET http://localhost:8080/clean-dummy-v2 \
  -H "Content-Type: application/json" \
  -i
```

## Issue Coupon Request V1

```bash
curl -X POST http://localhost:8080/api/v1/coupons/issue \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "policy_code": "BF-2025"
    }' \
    -i
```