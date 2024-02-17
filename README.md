# UserTransactions
Backend for money transfer application. Consists of two microservices Users and Transactions utilizing Postgres database, and Kafka and NATS for fast message exchange.
### Users
Users service can be used for new user registration and balance check. The service exposes two endpoints:
- _/createuser_ and
- _/balance_.
On user creation, users data are kept in database, and simultaneously published on Kafka. Published data is consumed by Transactions service.
Once the request for balance is sent, a NATS request is sent to Transactions service, which responds with actual account balance.
- _/createuser_
  - request
``` json
{
    "email": "john@doe.com"
}
```
  - response
```
confirmation message/error
```
- _/balance_
  - request
``` json
{
    "email": "john@doe.com",
}
```
  - response
``` json
{
    "email": "john@doe.com",
    "balance": 123
}
```

### Transactions
## Start
