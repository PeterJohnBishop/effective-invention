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
	// curl -X GET 'http://localhost:8080/'

	r.POST("auth/login", func(c *gin.Context) {
		pgdb.Login(db, c)
	})
	// curl -X POST 'http://localhost:8080/auth/login' \
	// 	-H 'Content-type: application/json' \
	// 	-d '{"email": "test2@gmail.com", "password": "MyPassword"}'

	r.POST("auth/register", func(c *gin.Context) {
		pgdb.RegisterUser(db, c)
	})
	// curl -X POST 'http://localhost:8080/auth/register' \
	// 	-H 'Content-type: application/json' \
	// 	-d '{"name": "test user", "email": "test@gmail.com", "password": "MyPassword"}'

	r.POST("auth/refresh", func(c *gin.Context) {
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
	// curl -X GET 'http://localhost:8080/api/users' \
	//  -H 'Content-Type: application/json' \
	//  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXJfMTQ3MTI3MjE2NzQzMzQ5MTQ5MTkiLCJleHAiOjE3Njc5MjYwMjAsImlhdCI6MTc2NzkyNTEyMH0.KihW-oX8wV_cHuPrgViL_A4RlLl9E_4pDxVmkKnRZnI' | jq

	r.GET("/users/:id", func(c *gin.Context) {
		pgdb.GetUserByID(db, c)
	})
	r.PUT("/users", func(c *gin.Context) {
		pgdb.UpdateUser(db, c)
	})
	r.PUT("/users/password", func(c *gin.Context) {
		pgdb.UpdatePassword(db, c)
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		pgdb.DeleteUserByID(db, c)
	})
	//	curl -X DELETE 'http://localhost:8080/api/users/user_14712721674334914919' \
	//	 -H 'Content-Type: application/json' \
	//	 -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXJfMTQ3MTI3MjE2NzQzMzQ5MTQ5MTkiLCJleHAiOjE3Njc5Mjg3NTQsImlhdCI6MTc2NzkyNzg1NH0.jaUjjSKununOhYTfAQjz8EBBJeMYiQfxJWfkQo6t7rU'
}
