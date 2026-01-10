package server

import (
	"database/sql"
	"effective-invention/server/cuapi"
	"effective-invention/server/pgdb"
	"effective-invention/server/pgdb/jwt"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeGin(db *sql.DB) {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	cuapi.ConnectClickUp()

	pgdb.CreateUsersTable(db)

	protected := r.Group("/api")
	protected.Use(jwt.JWTMiddleware())
	addOpenRoutes(r, db)
	addProtectedRoutes(protected, db)
	addClickUpRoutes(r)
	AddTestEndpoints(r)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	log.Println("Serving Postgres over Gin at localhost:8080.")
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
