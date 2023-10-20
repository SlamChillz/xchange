package db

import (
	"log"
	"os"
	"testing"
	"database/sql"

	"github.com/slamchillz/xchange/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// setup a new connection to the database
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
