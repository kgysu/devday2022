# SBB - DevDay 2022 - Go Webserver & API with PostgreSQL

This repo contains all the code shown in my DevDay-Talk.



## Links

Learn:

[learn go with tests](https://quii.gitbook.io/learn-go-with-tests/)

[Examples](https://github.com/kgysu/devday2022/tree/main/cmd/examples)

[Gophercon YouTube](https://www.youtube.com/c/GopherAcademy/playlists)



Go-Libs:

[Multiplexer Gorilla](https://github.com/gorilla/mux)

[PostgreSQL driver](https://github.com/jackc/pgx)

[Autocert](https://github.com/kgysu/devday2022/tree/main/cmd/examples/tls)

[JWT](https://github.com/golang-jwt/jwt)

[Logging](https://github.com/sirupsen/logrus)



## HowTo

Start at `cmd/starter/main.go`


### Generate SQLC Code

To generate the sqlc go-classes run:

```bash
cd pkg/database && sqlc generate
```


### Database migrations

To migrate a local postgresql database run:

```bash
# setup
migrate -source file://./migrations -database "postgresql://postgres:postgres@localhost:5432/devday?sslmode=disable" up 1

# teardown
migrate -source file://./migrations -database "postgresql://postgres:postgres@localhost:5432/devday?sslmode=disable" down 1
```

