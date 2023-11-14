postgres:
	@echo "Starting postgres..."
	@docker run --name goswap -p ${DBPORT}:${DBPORT} -e POSTGRES_USER=${DBUSER} -e POSTGRES_PASSWORD=${DBPASSWORD} -d postgres:12-alpine
	@echo "Postgres started."

createdb:
	@echo "Creating database..."
	@docker exec -it goswap createdb --username=${DBUSER} --owner=${DBUSER} ${DBNAME}
	@echo "Database created."

dropdb:
	@echo "Dropping database..."
	@docker exec -it goswap dropdb --username=${DBUSER} ${DBNAME}
	@echo "Database dropped."

migrateup:
	migrate -path db/migrations -database "${DBDRIVER}://${DBUSER}:${DBPASSWORD}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "${DBDRIVER}://${DBUSER}:${DBPASSWORD}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose down

mock:
	@echo "Generating mock..."
	@mockgen -package mockdb -destination db/mock/storage.go github.com/slamchillz/xchange/db/sqlc Store
	@echo "Mock generated."

sqlcinstall:
	@echo "Installing sqlc..."
	@docker pull sqlc/sqlc
	@echo "Sqlc installed."

sqlcinit:
	@echo "Initializing sqlc..."
	@sqlc init
	@echo "Sqlc initialized."

sqlcgenerate:
	@echo "Generating sqlc..."
	@sqlc generate
	@echo "Sqlc generated."

test:
	@echo "Testing..."
	@go test -v -cover ./...
	@echo "Tested."

server:
	go run main.go

proto:
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

evans:
	evans -r repl

.PHONY: createdb dropdb postgres migrateup migratedown sqlcinstall sqlcinit sqlcgenerate mock test server proto
