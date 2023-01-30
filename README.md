# Simple Wallet Api
* api start from port `:2565`
* use PostgreSQL
* database url get from environment variable name `DATABASE_URL`

### Url for test api
```console
https://wallet-kyxxckomzq-as.a.run.app/wallet
```

### Test Local
### Create and start containers
```console
docker-compose up
```console
### Stop and remove containers, networks
```console
docker-compose down --rmi local -v
```

### To run integration test
```console
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests
```
### To tear down integration test
```console
docker-compose -f docker-compose.test.yml down --rmi local -v
```

* Postman collection สำหรับทดสอบ API ทั้งหมด
	- สำหรับทดสอบบน localhost [postman collection](wallet-cloud.postman_collection.json)
	- สำหรับทดสอบบน localhost [postman collection](localhost-wallet.postman_collection.json)

#### Technical Details: List all wallets
* GET /wallet
* Response Body
```json
[
    {
        "wallet_id": "1",
        "balance": 1000,
        "status": "Active",
        "created_at": "2023-01-27T12:30:00Z"
    },
    {
        "wallet_id": "2",
        "balance": 2000,
        "status": "Active",
        "created_at": "2023-01-27T12:30:00Z"
    }
]
```

#### Technical Details: Create a new wallet
* POST /wallet
* Request Body
```json
{
	"balance": 1000
}
```
* Response Body
```json
{
	"wallet_id": "1",
	"balance": 1000,
	"status": "Active",
	"created_at": "2023-01-27T12:30:00Z"
}
```

#### Technical Details: Get a wallet's detail 
* GET /wallet/:id
* :id = 1
* Response Body
```json
{
    "wallet_id": "1",
    "balance": 1000,
    "status": "Active",
    "created_at": "2023-01-27T12:30:00Z"
}
```

#### Technical Details: Add balance to a wallet
* PUT /wallet/:id
* :id = 1
* Request Body
```json
{
	"balance": 1000,
    "operation": "Add"
}
```
* Response Body
```json
{
	"wallet_id": "1",
	"balance": 2000,
	"status": "Active",
	"created_at": "2023-01-27T12:30:00Z"
}
```

#### Technical Details: Deduct balance from a wallet
* PUT /wallet/:id
* :id = 1
* Request Body
```json
{
	"balance": 500,
    "operation": "Deduct"
}
```
* Response Body
```json
{
	"wallet_id": "1",
	"balance": 500,
	"status": "Active",
	"created_at": "2023-01-27T12:30:00Z"
}
```

#### Technical Details: Deactivate & Activate a wallet
* PUT /wallet/:id/status
* :id = 1
* Request Body
```json
{
    "operation": "Active"
}
```
* Response Body
```json
{
	"wallet_id": "1",
	"balance": 500,
	"status": "Active",
	"created_at": "2023-01-27T12:30:00Z"
}
```

* Request Body
```json
{
    "operation": "Deactive"
}
```
* Response Body
```json
{
	"wallet_id": "1",
	"balance": 500,
	"status": "Deactive",
	"created_at": "2023-01-27T12:30:00Z" 
}
```
