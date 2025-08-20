## SignUp

```bash
curl -X POST "http://localhost:8081/api/v1/auth/sign-up" \
     -H "Content-Type: application/json" \
     -d '{
           "fullName": "John Doe",
           "email": "jhondoe@example.com",
           "password": "mypassword"
         }'
```

## SignIn

```bash
curl -X POST "http://localhost:8081/api/v1/auth/sign-in" \
     -H "Content-Type: application/json" \
     -d '{
           "email": "jhondoe@example.com",
           "password": "mypassword"
         }'
```

## GetProfileMe

```bash
curl -X GET "http://localhost:8081/api/v1/users/profiles/me" \
     -H "Authorization: Bearer TOKEN" \
     -H "Accept: application/json"
```

## SearchTickets

```bash
curl -X GET "http://localhost:8081/api/v1/tickets" \
     -H "Authorization: Bearer TOKEN" \
     -H "Accept: application/json"
     
curl -X GET "http://localhost:8081/api/v1/tickets?search=1" \
     -H "Authorization: Bearer TOKEN" \
     -H "Accept: application/json"
     
curl -X GET "http://localhost:8081/api/v1/tickets?status=OPEN&priority=HIGH" \
     -H "Accept: application/json"
     
curl -X GET "http://localhost:8081/api/v1/tickets?date=TODAY" \
     -H "Authorization: Bearer TOKEN" \
     -H "Accept: application/json"
     
curl -X GET "http://localhost:8081/api/v1/tickets?search=cable&status=RESOLVED&priority=LOW&date=YEAR" \
     -H "Authorization: Bearer TOKEN" \
     -H "Accept: application/json"
```

## CreateTicket

## UpdateTicket

## CreateTicketReply

## Dashboards

```bash
curl -X GET "http://localhost:8081/api/v1/dashboards" \
     -H "Authorization: Bearer TOKEN" \
     -H "Accept: application/json"
```