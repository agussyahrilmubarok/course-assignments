# Get All Categories

```bash
curl -X GET "http://localhost:8080/api/categories" \
     -H "Accept: application/json"
```

# Get Category by ID

```bash
curl -X GET "http://localhost:8080/api/categories/1" \
     -H "Accept: application/json"
```

# Create New Category (admin only)

```bash
curl -X POST "http://localhost:8080/api/categories" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your_admin_jwt_token>" \
     -d '{"name": "New Category"}'
```

# Update Category by ID (admin only)

```bash
curl -X PUT "http://localhost:8080/api/categories/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your_admin_jwt_token>" \
     -d '{"name": "Updated Category Name"}'
```

# Delete Category by ID (admin only)

```bash
curl -X DELETE "http://localhost:8080/api/categories/1" \
     -H "Authorization: Bearer <your_admin_jwt_token>"
```