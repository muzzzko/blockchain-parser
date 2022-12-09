

## Tools
- mockgen 

to install run
```
go install github.com/golang/mock/mockgen@v1.6.0
```
- docker

## Start

1) Clone repository to github/blockchain-parser folder

```
git clone git@github.com:muzzzko/blockchain-parser.git ./github/blockchain-parser/
cd github/blockchain-parser/
```

2) Build image
```
make build-service-image
```

3) Start server
```
make run-in-docker
```

to stop server 
```
make stop-in-docker
```

## Run tests

```
make run-test-in-docker
```

## Settings

- BLOCKCHAIN_PARSER_PARSER_WORKER_COUNT_WORKERS - sets count workers which will parse block (required)
- BLOCKCHAIN_PARSER_PARSER_WORKER_INTERVAL - sets waiting interval for workers (required) (time.Duration format)
- BLOCKCHAIN_PARSER_PARSER_WORKER_START_BLOCK_NUMBER - sets initial block (default: -1). If it was set, then workers start parsing from particular block. If it wasn't set, then start from the last block. (hex format)
- BLOCKCHAIN_PARSER_PARSER_WORKER_PREDEFINED_ADDRESSES - sets initial addresses for subscribing. It's useful in case when you don't want to make a transaction but you need to check GetTransactions method. 
Set BLOCKCHAIN_PARSER_PARSER_WORKER_START_BLOCK_NUMBER and set BLOCKCHAIN_PARSER_PARSER_WORKER_PREDEFINED_ADDRESSES with addresses from this block. Addresses should split with ',' 
```
Example: BLOCKCHAIN_PARSER_PARSER_WORKER_PREDEFINED_ADDRESSES=0xa855d1198c67839e596b9a5d7c46f8ea31cfefde,0xfd4492e70df97a6155c6d244f5ec5b5a39b6f096
```

## Improvements
1) In current implementation only one instance can work, but you can set several workers inside this instance to parallel parsing. 
For run several instances you need to replace in memory store to DB (e.g. postgres)
2) If you need to get all user transactions, you will have to make some changes. Service should pull all transactions from the beginning and stor them to DB. 
Service will do it once during first start. After that service starts pulling from last parsed block and stores all transaction instead of storing only transactions with subscribed addresses.  
To store all transaction (three fields: From, To, Value) we need 2.4TB 
3) Need to implement state machine for block processing
4) Add err codes to response
5) Subscriber must be different service, now it's inside parser
6) Parser must return errors
7) Replace default logger to zap for instance. I used default to fit task requirements

## Notes
I suppose that skipping external packages increases security. I afforded myself to use mockgen and monkey because they are used only for testing and won't be inside production build 
