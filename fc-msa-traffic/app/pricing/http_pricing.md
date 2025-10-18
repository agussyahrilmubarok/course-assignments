# HTTP PRICING SERVICE

## Create Pricing Rule

`api/v1/pricings/rules`

```bash
curl -X POST http://localhost:8084/api/v1/pricings/rules \
    -H "Content-Type: application/json" \
    -d '{
        "product_id": "",
        "default_markup": 0.15,
        "default_discount": 0.05,
        "stock_threshold": 50,
        "markup_increase": 0.10,
        "discount_reduction": 0.02
    }' \
    -i
```

## Get Pricing

`api/v1/pricings/:id`

```bash
curl -X GET http://localhost:8084/api/v1/pricings/beb65fe0-77df-4815-967f-2294f762d8c4 \
    -H "Content-Type: application/json" \
    -i
```
