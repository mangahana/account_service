include .env
export

run:
	go run ./cmd/main.go

database_up:
	docker run --name test-db --rm \
	-e POSTGRES_USER=test \
	-e POSTGRES_PASSWORD=test \
	-e POSTGRES_DB=test \
	-e PGDATA=/var/lib/postgresql/data \
	-p 5432:5432 \
	-v ${CURDIR}/migrations:/docker-entrypoint-initdb.d \
	-d postgres:15.3-bullseye
	
database_down:
	docker stop test-db

proto_compile:
	protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/account.proto