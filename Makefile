DB_SOURCE ?= "postgresql://root:Fq9zkLWA2ZBAhq@localhost:5432/simple_payment?sslmode=disable"

postgres:
	docker run --name postgres18 --network payment-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Fq9zkLWA2ZBAhq -d postgres:18

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root simple_payment

dropdb:
	docker exec -it postgres18 dropdb simple_payment

migrateup:
	migrate -path db/migration -database $(DB_SOURCE) -verbose up

migrateup1:
	migrate -path db/migration -database $(DB_SOURCE) -verbose up 1


migratedown:
	migrate -path db/migration -database $(DB_SOURCE) -verbose down

migratedown1:
	migrate -path db/migration -database $(DB_SOURCE) -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/ryannguyen1105/Simplepayment/db/sqlc Store

.PHONY: postgres	createdb	dropdb	migrateup	migratedown	migrateup1	migratedown1	sqlc	test	server	mock	db_docs	db_schema