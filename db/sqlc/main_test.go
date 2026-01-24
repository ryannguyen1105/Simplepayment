package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ryannguyen1105/Simplepayment/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:Fq9zkLWA2ZBAhq@localhost:5432/simple_payment?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
