package main

import (
	"database/sql"
	"effective-invention/server"
	"effective-invention/server/pgdb"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	fmt.Println("[ effective-invention ]( launching )")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db = pgdb.PostGres()
	defer db.Close()
	server.ServeGin(db)
}
