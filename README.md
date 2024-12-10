## Project specific folders

- cmd
    - contains the entrypoint of the application, which is the `main.go` file.
- config
    - contains the yml config files for each environment and `config.go`, which contains the implementation
      of the schema, and it's interpretation.
- internal
    - contains all the business logic inspired
    - in this folder and sub-folders, you will only find domain models, domain services or the port interfaces.
- db
    - contains all db migrations.

### How to write a migration

* add your migration under db/migrations
* Migrations will be run automatically when app start creating the store


### How to work
For develop you need to run `docker-compose up -d --build` and then go to container and run go run `go run cmd/api/main.go`
after these, database migration will run, and you can start to send request to application
also if you want you can uncomment the `entrypoint: ["/usr/local/go/bin/go", "run", "cmd/api/main.go"]` in docker-compose file 
and run `docker-compose up app`

For production you need to `docker build --target main  -t bank-api  -f docker/Dockerfile  .  --no-cache` and then run 
docker file and make a request to it.




## Routes
### USER
**user/create**
   POST create user with account type of normal and value 0 w
`{
"first_name": "mozhgan",
"last_name":"parsa",
"email":"mozhgan@gmail.com",
"active":true
}
`

**user/find**
    GET get user
`    http://localhost:8080/v1/user/find?id=16`

### Transfer
`v1/transfer/money`

{
"from": 1,
"to": 4,
"amount": 300
}

### Account
create account
`v1/account/create`

{
"user_id": 1,
"type" : "normal",
"balance": 100
}

update account
`/v1/account/update`

{
"user_id": 1,
"balance": 7777
}
