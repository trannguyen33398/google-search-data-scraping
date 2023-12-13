package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

func InitConnection() (*sql.DB, error) {


	// Build the connection string
	
	db, err := sql.Open("postgres","postgres://admin:admin@localhost/test?sslmode=disable")

	if err != nil {
		fmt.Println("Postgres", "InitPostgresConnection()", err)
		log.Panic()
	}

	fmt.Println("Connect Postgres successfully")

	return db, err
}