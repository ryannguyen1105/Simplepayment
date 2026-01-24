package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ryannguyen1105/Simplepayment/api"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
	"github.com/ryannguyen1105/Simplepayment/util"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:Fq9zkLWA2ZBAhq@localhost:5432/simple_payment?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start sever:", err)
	}
}
