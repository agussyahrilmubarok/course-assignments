# SignUp

`POST /api/v1/auth/sign-up`

```bash
curl -X POST http://localhost:8081/api/v1/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Agus Syahril",
    "email": "agus@example.com",
    "password": "secret123"
  }' \
  -i
```

# SignIn

`POST /api/v1/auth/sign-in`

```bash
curl -X POST http://localhost:8081/api/v1/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "agus@example.com",
    "password": "secret123"
  }' \
  -i
```

# Validate Token

`POST /api/v1/auth/validate-token`

```bash
curl -X POST http://localhost:8081/api/v1/auth/validate-token \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..."
  }' \
  -i
```

# Get User

`GET /api/v1/users/me`

```bash
curl -X GET "http://localhost:8081/api/v1/users/me" \
  -H "Content-Type: application/json" \
  -H "X-USER-ID: a6c3dfb0-80e8-43f2-b62e-a36ecf289e31" \
  -i
```