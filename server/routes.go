package server

import "github.com/gin-gonic/gin"

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
