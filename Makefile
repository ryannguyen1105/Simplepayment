postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Fq9zkLWA2ZBAhq -d postgres:18

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root simple_payment

dropdb:
	docker exec -it postgres18 dropdb simple_payment

migrateup:
	migrate -path db/migration -database "postgresql://root:Fq9zkLWA2ZBAhq@localhost:5432/simple_payment?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:Fq9zkLWA2ZBAhq@localhost:5432/simple_payment?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres	createdb	dropdb	migrateup	migratedown	sqlc	test