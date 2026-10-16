set dotenv-load := true

build:
	go build cmd/api/main.go

run:
	go run cmd/api/main.go

db-up:
	goose -dir migrations turso "$LIBSQL_URL?authToken=$LIBSQL_AUTH_TOKEN" up

db-status:
	goose -dir migrations turso "$LIBSQL_URL?authToken=$LIBSQL_AUTH_TOKEN" status
