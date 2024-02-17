# UserTransactions
Backend for money transfer application. Consists of two microservices Users and Transactions utilizing Postgres database, and Kafka and NATS for fast message exchange.
### Users
Users service can be used for new user registration and balance check. The service exposes two endpoints:
- _/createuser_ and
- _/balance_.  
On user creation, users data are kept in database, and simultaneously published on **Kafka**. Published data is consumed by Transactions service.
Once the request for balance is sent, a **NATS** request is sent to Transactions service, which responds with actual account balance.
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
Transactions service keeps track of commited transactions and user's account balance. The service exposes two endpoints  
- _/deposit_ and
- _transfer_.  
Once the new user is created, Transactions service is notified over Kafka. The message is processed by **group consumer**, and based on the message, new account is created with initial balance of 0. Transactions service also takes part at responding to balance queries: the service receives account id, executes database lookup, and responds with curent balance. When transfer request is created, the service check if transfer can be valid, and only than proceeds to withdraw and deposit funds from one account to another. Record of each transaction is keept in database.
- _/deposit_
  - request
``` json
{
    "userId": 24,
    "amount": 11111
}
```
  - response
```
{
    "balance": 123
}
```
- _/transfer_
  - request
``` json
{
    "from": 2,
    "to": 1,
    "amount": 111
}
```
  - response
```
confirmation message/error
```
## Start
`cd` to projet root direcotry. There ***docker-compose.yml*** file can be found. Use `docker-compose up -d` to start it.  
Once it is up and runnig Postgres (:5432), Kafka (:9092), Zookeeper (:2181), and NATS (:4222) will be available.  
To start Users service execute `go run ./cmd/usersd/ .`  
To start Transactions service execute `go run ./cmd/transactionsd/ .`
Both services are by default configured to connect to Postgres, Kafka, and NATS.  
For detailed configuration both services expose folowing flags:
- service.port - Port for service to run on
- db.host - Database Host
- db.port - Database Port
- db.name - Database Name
- db.user - Database User
- db.password - Database Password
- nats.subject - NATS subject pattern
- kafka.brokers - Kafka broker address
- kafka.topic - Kafka topic name
- kafka.topic - Kafka topic name (***only for transactions service***)
  Default flag values can be found at the top of main.go file for each service.
