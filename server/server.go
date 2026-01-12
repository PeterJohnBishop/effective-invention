package server

import (
	"effective-invention/server/amazonwebservices"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"
)

var aws_client aws.Config

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

	aws_config := amazonwebservices.StartAws()
	s3_client := amazonwebservices.ConnectS3(aws_config)
	log.Println(s3_client)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	addS3Routes(s3_client, r)

	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	config := fmt.Sprintf(":%s", port)
	log.Printf("Serving Gin ats %s:%s", baseUrl, port)
	r.Run(config)
}
