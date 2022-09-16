# SBB - DevDay 2022 - Go Webserver & API with PostgreSQL

This repo contains all the code shown in my DevDay-Talk.


## Generate SQLC Code

To generate the sqlc go-classes run:

```bash
cd pkg/database && sqlc generate
```


## Database migrations

To migrate a local postgresql database run:

```bash
# setup
migrate -source file://./migrations -database "postgresql://postgres:postgres@localhost:5432/devday?sslmode=disable" up 1

# teardown
migrate -source file://./migrations -database "postgresql://postgres:postgres@localhost:5432/devday?sslmode=disable" down 1
```



## Links

Useful libraries:

[Multiplexer Gorilla](https://github.com/gorilla/mux)

[Logging](https://github.com/sirupsen/logrus)

[PostgreSQL driver](https://github.com/jackc/pgx)
