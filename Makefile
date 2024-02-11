include .env
export

run:
	go run ./cmd/auth/main.go

build:
	go build ./.bin/auth-ms ./cmd/auth


migrate_up:
	migrate -path ${PATH_TO_MIGRATIONS} -database postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=${PGSSLMODE} up


