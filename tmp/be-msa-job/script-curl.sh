#!/bin/sh

## Init accounts, cards, loans
curl -X GET http://localhost:8072/examplebank/accounts/seeds/init
curl -X GET http://localhost:8072/examplebank/cards/seeds/init
curl -X GET http://localhost:8072/examplebank/loans/seeds/init

# SCHEMA MANUAL TEST
# 1. Should Be Success
## Fetch Customer Details
curl -X GET "http://localhost:8072/examplebank/accounts/api/customers/fetchCustomerDetails?mobileNumber=0812345678" \
     -H "Accept: application/json"

# SCHEMA MANUAL TEST
# 2. Should Be Null On Cards Service
## Shutdown Cards Service
## Fetch Customer Details
curl -X GET "http://localhost:8072/examplebank/accounts/api/customers/fetchCustomerDetails?mobileNumber=0812345678" \
     -H "Accept: application/json"

# SCHEMA MANUAL TEST
# 3. Should Be Null On Loans Service
## Shutdown Loans Service
## Fetch Customer Details
curl -X GET "http://localhost:8072/examplebank/accounts/api/customers/fetchCustomerDetails?mobileNumber=0812345678" \
     -H "Accept: application/json"

# SCHEMA MANUAL TEST
# 3. Should Be Got Limit
## Loans Service
## Fetch Cards Details
for i in {1..100}; do
  curl -s -o /dev/null -w "Request $i: %{http_code}\n" \
  -X GET "http://localhost:8072/examplebank/cards/api/fetch?mobileNumber=0812345678" \
  -H "Accept: application/json" \
  -H "user: potatohuman" &
done
wait
## AND will see Request 429: Too Many Requests.
## docker exec -it redis redis-cli KEYS *

# Create A New Account Should Use Token
## Use Postman for easy test

# Clean accounts, cards, loans
curl -X GET http://localhost:8072/examplebank/accounts/seeds/clean
curl -X GET http://localhost:8072/examplebank/cards/seeds/clean
curl -X GET http://localhost:8072/examplebank/loans/seeds/clean


