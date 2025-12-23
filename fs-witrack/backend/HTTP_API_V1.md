# HTTP API V1

## Sign Up

```bash
curl -X POST http://localhost:8081/api/v1/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

## Sign In

```bash
curl -X POST http://localhost:8081/api/v1/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

## Get Me

```bash
curl -X GET http://localhost:8081/api/v1/users/profiles/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token"
```

## Create Ticket

```bash
curl -X POST http://localhost:8081/api/v1/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token" \
  -d '{
    "title": "Create Ticket 1 Title",
    "description": "Create Ticket 1 Description",
    "status": "OPEN",
    "priority": "HIGH"
  }'
```

```bash
curl -X POST http://localhost:8081/api/v1/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token" \
  -d '{
    "description": "Create Ticket 1 Description",
    "status": "OPEN",
    "priority": "HIGH"
  }'
```

## Update Ticket

```bash
curl -X PUT http://localhost:8081/api/v1/tickets/34b28f8c-2e2c-483b-9955-fe367476e9bc \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token" \
  -d '{
    "title": "Create Ticket 1 Title Update",
    "description": "Create Ticket 1 Description",
    "status": "OPEN",
    "priority": "HIGH"
  }'
```

## Delete Ticket

```bash
curl -X DELETE http://localhost:8081/api/v1/tickets/ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer token"
```