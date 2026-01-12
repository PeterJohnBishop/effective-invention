package server

import (
	"effective-invention/server/github"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func ServeGin() {
	log.Println("Orderings Gin")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	github.GetGit()

	r.GET("/ping", func(c *gin.Context) {
		github.GetUserRepos("PeterJohnBishop")
		c.String(200, "pong")
	})
	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	config := fmt.Sprintf(":%s", port)
	log.Printf("Serving Gin ats %s:%s", baseUrl, port)
	r.Run(config)
}
