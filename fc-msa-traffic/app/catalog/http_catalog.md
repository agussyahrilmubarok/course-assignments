# HTTP CATALOG SERVICE

## Get Catalog Products

`api/v1/catalogs/products`

```bash
curl -X GET http://localhost:8082/api/v1/catalogs/products \
    -H "Content-Type: application/json" \
    -i
```

## Get Catalog Product

`api/v1/catalogs/products/id`

```bash
curl -X GET http://localhost:8082/api/v1/catalogs/products/ \
    -H "Content-Type: application/json" \
    -i
```

## Reverse Stock Product

`api/v1/catalogs/products/reverse`

```bash
curl -X POST http://localhost:8082/api/v1/catalogs/products/reverse \
    -H "Content-Type: application/json" \
    -d '{
        "product_id": "",
        "quantity": 1
    }'
    -i
```

## Release Stock Product

`api/v1/catalogs/products/release`

```bash
curl -X POST http://localhost:8082/api/v1/catalogs/products/release \
    -H "Content-Type: application/json" \
    -d '{
        "product_id": "",
        "quantity": 1
    }'
    -i
```
