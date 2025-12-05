# BackerHub Backend

## Run BackerHub Backend

```bash
# Run in development
make compose/postgres/up
make backend/dev
make seeder/run
make compose/postgres/down

# Run in stage
make compose/stage/up
make compose/stage/down
```

## Technology Stack

* Golang
* Gin
* Gin Templates
* GORM
* Zap
* Postgres

## References

