package server

import (
	"database/sql"
	"effective-invention/server/pgdb"
	"log"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func ServeGin() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	db = pgdb.PostGres()

	AddTestEndpoints(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	log.Println("Serving Postgres over Gin at localhost:8080.")
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
