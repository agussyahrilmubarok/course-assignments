# HTTP ORDER SERVICE

## Create Order

`api/v1/orders`

```bash
curl -X POST http://localhost:8083/api/v1/orders/flash \
    -H "Content-Type: application/json" \
    -d '{
        "order_items": [
            {
                "product_id": "1e34b1e7-7ba1-4633-b5db-d3d77efb3d38",
                "quantity": 5
            }
        ],
        "user_id": "user-x-1"
    }' \
    -i
```

## Cancel Order

`api/v1/orders`

```bash
curl -X POST http://localhost:8083/api/v1/orders/cancel \
    -H "Content-Type: application/json" \
    -d '{
        "order_id": ""
        "user_id": "user-x-1"
    }' \
    -i
```

## Get Order

`api/v1/orders`

```bash
curl -X GET http://localhost:8083/api/v1/orders/ \
    -H "Content-Type: application/json" \
    -i
```

## Example pricing rules

```bash
{
  "default_discount": 0.05,
  "default_markup": 0.15,
  "discount_reduction": 0.02,
  "markup_increase": 0.03,
  "product_id": "1e34b1e7-7ba1-4633-b5db-d3d77efb3d38",
  "stock_threshold": 50
}
```
