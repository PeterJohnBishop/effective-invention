package server

import (
	"database/sql"
	"effective-invention/server/pgdb"

	"github.com/gin-gonic/gin"
)

func AddTestEndpoints(r *gin.Engine) {

	r.POST("/raw", HandleRawData())
	// curl -X POST 'http://localhost:8080/raw' \
	//  -H "Content-Type: application/json" \
	//  -d '{"title": "My New Post", "body": "This is the content.", "userId": 1}'

	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})
	// curl http://localhost:8080/health

	r.POST("/build", HandleTypeMap())
	// curl -X POST 'http://localhost:8080/build?name=Post' \
	//  -H "Content-Type: application/json" \
	//  -d '{"title": "My New Post", "body": "This is the content.", "userId": 1}'

}

func addOpenRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"repo": "https://github.com/PeterJohnBishop/effective-invention",
		})
	})
	r.POST("auth/login", func(c *gin.Context) {
		pgdb.Login(db, c)
	})
	r.POST("auth/register", func(c *gin.Context) {
		pgdb.RegisterUser(db, c)
	})
	r.GET("auth/refresh", func(c *gin.Context) {
		pgdb.Refresh(c)
	})
	r.POST("auth/logout", func(c *gin.Context) {
		pgdb.Logout(c)
	})
}

func addProtectedRoutes(r *gin.RouterGroup, db *sql.DB) {

	r.GET("/users", func(c *gin.Context) {
		pgdb.GetUsers(db, c)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		pgdb.GetUserByID(db, c)
	})
	r.PUT("/users", func(c *gin.Context) {
		pgdb.UpdateUser(db, c)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		pgdb.DeleteUserByID(db, c)
	})
}
