package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ryannguyen1105/Simplepayment/api"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:Fq9zkLWA2ZBAhq@localhost:5432/simple_payment?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start sever:", err)
	}
}
