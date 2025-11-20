api/dev:
	go run cmd/api/main.go --config config.yml

#####################################################################################
### migrater
#####################################################################################
migrater/up:
	go run cmd/migrater/main.go \
	--db postgres://postgres:my-secret-pw@localhost:5432/coupon_db?sslmode=disable \
	--dir migrations \
	--action up

migrater/down:
	go run cmd/migrater/main.go \
	--db postgres://postgres:my-secret-pw@localhost:5432/coupon_db?sslmode=disable \
	--action down

# go run cmd/migrater/main.go --db <URL> --action up
# go run cmd/migrater/main.go --db <URL> --action down
# go run cmd/migrater/main.go --db <URL> --action drop
# go run cmd/migrater/main.go --db <URL> --action version
# go run cmd/migrater/main.go --db <URL> --action force --force 3
# go run cmd/migrater/main.go \
#   --db postgres://postgres:my-secret-pw@localhost:5432/coupon_db?sslmode=disable \
#   --dir migrations \
#   --action up
migrater/command:
	go run cmd/migrater/main.go \
	--db $(DB_URL) \
	--dir $(DIR) \
	--action $(C)


