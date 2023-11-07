package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/slamchillz/xchange/api"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
)

func main() {
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	server, err := api.NewServer(config, db.NewStorage(conn))
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	log.Fatal(server.Start(":8080"))
}
