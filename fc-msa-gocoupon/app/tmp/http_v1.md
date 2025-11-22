# HTTP v1 Example

## Init Dummy Coupon Policy V1

```bash
curl -X GET http://localhost:8080/init-dummy-v1 \
  -H "Content-Type: application/json" \
  -i
```

## Clean Dummy Coupon Policy V1

```bash
curl -X GET http://localhost:8080/clean-dummy-v1 \
  -H "Content-Type: application/json" \
  -i
```

## Issue Coupon Request V1

```bash
curl -X POST http://localhost:8080/api/v1/coupons/issue \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "policy_code": "C-V1-1"
  }' \
  -i
```

## Issue Coupon Request V1 Loop (11 requests)

```bash
for i in {1..11}
do
  USER_ID="USER_$i"
  echo "Request ke-$i dengan X-USER-ID: $USER_ID"

  curl -X POST http://localhost:8080/api/v1/coupons/issue \
    -H "Content-Type: application/json" \
    -H "X-USER-ID: $USER_ID" \
    -d '{"policy_code": "C-V1-1"}' \
    -i

  echo "" 
done
```

## Use Coupon Request V1

```bash
curl -X POST http://localhost:8080/api/v1/coupons/use \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "coupon_code": "",
    "order_id": "ORDER-12345"
  }' \
  -i
```

## Cancel Coupon Request V1

```bash
curl -X POST http://localhost:8080/api/v1/coupons/cancel \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -d '{
    "coupon_code": ""
  }' \
  -i
```

## Find Coupon By Code

```bash
curl -X GET http://localhost:8080/api/v1/coupons/78e6b21b-2f98-4fa5-bff7-2cb4ba96db61 \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: USER_1" \
  -i
```
