#### Create a migration

```bash
$ go run -tags sqlite3 github.com/golang-migrate/migrate/v4/cmd/migrate create -dir=db/migrations -ext=.sql <migration_name>
```

#### parse configs and secrets, apply migrations, and compile queries before running the program

```bash
go generate && go run .
```
