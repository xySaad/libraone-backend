### Developing

#### Env

GITEA_CLIENT_SECRET is a secret for a gitea application with scopes: `read:user`

#### Create a migration

```bash
$ go run -tags sqlite3 github.com/golang-migrate/migrate/v4/cmd/migrate create -dir=db/migrations -ext=.sql <migration_name>
```

#### parse configs and secrets, apply migrations, and compile queries before running the program

```bash
go generate && go run .
```

An example of using this API can be found at https://github.com/xySaad/libraone-frontend
