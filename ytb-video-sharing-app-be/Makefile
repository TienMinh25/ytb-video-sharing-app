MYSQL_ROOT_PASSWORD:=admin123
MYSQL_DATABASE:=ytb
MYSQL_HOST:=localhost
MYSQL_PORT:=3307
MYSQL_URL := "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)"

MIGRATE_DIR :=db/migrations

migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

migrate-up:
	@migrate -database $(MYSQL_URL) -path db/migrations up

migrate-down:
	@migrate -database $(MYSQL_URL) -path db/migrations down

fix-dirty-db:
	@migrate -database $(MYSQL_URL) -path db/migrations force $(version)

generate-public-key: generate-private-key
	@openssl rsa -pubout -in jwtRSA256.key -out jwtRSA256.key.pub
	
generate-private-key:
	@openssl genpkey -algorithm RSA -out jwtRSA256.key

swagger-generate:
	@swag init -g cmd/main.go

swagger-format:
	@swag fmt

generate-mock:
	@go generate ./...

tests-run:
	@go test -v -count=1 ./... 2>&1 | grep -v "no test files"

tests-cover:
	@go test -coverprofile=internal/repository/test-cover.out -count=1 -v ./internal/repository/
	@go tool cover -html=internal/repository/test-cover.out

tests-clear:
	@rm internal/repository/test-cover.out
