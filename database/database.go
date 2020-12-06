package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DbConn *sql.DB

func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:729184Inna_@tcp(127.0.0.1:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
		return
	}
}
