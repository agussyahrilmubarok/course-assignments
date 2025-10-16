# HTTP ACCOUNT SERVICE

## Sign Up

`api/v1/accounts/sign-up`

```bash
curl -X POST http://localhost:8081/api/v1/accounts/sign-up \
    -H "Content-Type: application/json" \
    -d '{
        "name": "John Doe",
        "email": "johndoe@mail.com",
        "password": "P@ssw0rd"
    }' \
    -i
```

## Sign In

`api/v1/accounts/sign-in`

```bash
curl -X POST http://localhost:8081/api/v1/accounts/sign-in \
    -H "Content-Type: application/json" \
    -d '{
        "email": "johndoe@mail.com",
        "password": "P@ssw0rd"
    }' \
    -i
```

## Validate

`api/v1/accounts/validate`

```bash
curl -X POST http://localhost:8081/api/v1/accounts/validate \
    -H "Content-Type: application/json" \
    -d '{
        "token": ""
    }' \
    -i
```

## Get Me

`api/v1/accounts/me`

```bash
curl -X GET http://localhost:8081/api/v1/accounts/me \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjA2OTYyMzQsImlhdCI6MTc2MDYwOTgzNCwidXNlcl9pZCI6ImRmNzgwMWM2LWM0ZDUtNDI2YS1iYTM2LTEyMGQxZGIxODk1MiJ9.0qzNTRYcDKX5TbETPPFWpJ_zke-GiyHod66ouPm9G3c" \
    -i
```
