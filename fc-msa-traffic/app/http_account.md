# Sign Up

`api/v1/account/sign-up`

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
