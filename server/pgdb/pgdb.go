package pgdb

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func PostGres() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening db: ", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the db: ", err)
		return nil
	}

	log.Println("Connected to Postgres volume in Docker!")
	return db

}
