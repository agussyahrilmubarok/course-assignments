# Get Current User Profile (User / Admin)

```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <your_jwt_token>"
```

# Update Current User Profile

```bash
curl -X PUT http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name",
    "email": "newemail@example.com"
  }'
```

# Get All Users (admin only)

```bash
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer <admin_jwt_token>"
```

# Delete User by ID (admin only)

```bash
curl -X DELETE http://localhost:8080/api/users/3 \
  -H "Authorization: Bearer <admin_jwt_token>"
```