package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	sqlc "github.com/andreaswiidi/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

var testQueries *sqlc.Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal()
	}

	testQueries = sqlc.New(testDB)

	os.Exit(m.Run())
}
