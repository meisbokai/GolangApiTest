# GolangApiTest

GolangApiTestis an API server that interacts with a Postgresql database to manage a database of Users

## Installation
Clone this repo into your desired directory

```bash
git clone https://github.com/meisbokai/GolangApiTest
```

## Usage
Execution of this repository assumes that the device running it has the following:
 - Go (Implemented on v1.21.5)
 - Postgres SQL 

### Dependency Installation
```bash
go get -d ./...
```

### Setting up the database
We will need to create the table that will store the user data
```bash
make mig-up
```
For the sake of testing, `user`s can be seeded into the database
```bash
make seed
```

### Running the server
Use the following command to run the server
```bash
make serve
```

### End point documentation

API documentation is available on swagger in the following url:
http://localhost:3000/api/v1/swagger/index.html

Note: url is only active when the server is running

For ease of usage, the jwt token found below will grant admin access to all the endpoints. It can also be obtained using the login endpoint with the credentials found below.

* Admin Credentials
    ```json
    {
    "email": "test1@example.com",
    "password": "12345"
    }
    ```
* Login Result
    ```json
    {
        "id": "6946bc07-558b-47e7-bb47-f80e50e15e3c",
        "username": "test user 1",
        "email": "test1@example.com",
        "password": "$2a$10$aQ2YHu20mmm.K7vi2qB.RufPG0i8VVlaBcnatoIYxlDUZPlBauhRG",
        "role_id": 1,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2OTQ2YmMwNy01NThiLTQ3ZTctYmI0Ny1mODBlNTBlMTVlM2MiLCJVc2VybmFtZSI6InRlc3QgdXNlciAxIiwiSXNBZG1pbiI6dHJ1ZSwiRW1haWwiOiJ0ZXN0MUBleGFtcGxlLmNvbSIsIlBhc3N3b3JkIjoiJDJhJDEwJGFRMllIdTIwbW1tLks3dmkycUIuUnVmUEcwaThWVmxhQmNuYXRvSVl4bERVWlBsQmF1aFJHIiwiaXNzIjoibWVpc2Jva2FpIiwiZXhwIjoxNzA0NjUyNjE4LCJpYXQiOjE3MDQ2MzQ2MTh9.NTSuDp3vz22QFEB_bueFS5gARhG6xlrgJtKzAFanB2Y",
        "created_at": "2024-01-07T21:36:45.054134+08:00",
        "updated_at": null
    }
    ```

To use the jwt token, append the token string to `jwt <token>` and insert it to an `Authorization` header. Examples can be found in the swagger endpoint documentation. 


## Testing
A coverage test can be conducted using the following command
```
make coverage
```


