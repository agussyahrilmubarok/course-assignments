# Sign Up

`/api/v1/members/sign-up`

```bash
curl -X POST http://localhost:8081/api/v1/members/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secret123"
  }' \
  -i
```

# Sign In

`/api/v1/members/sign-in`

```bash
curl -X POST http://localhost:8081/api/v1/members/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secret123"
  }' \
  -i
```

# Validate Token

`/api/v1/members/validate`

```bash
curl -X POST http://localhost:8081/api/v1/members/validate \
  -H "Content-Type: application/json" \
  -d '{
    "token": ""
  }' \
  -i
```