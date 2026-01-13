package server

import (
	"effective-invention/server/amazonwebservices"
	"effective-invention/server/amazonwebservices/database"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"
)

var aws_client aws.Config

func ServeGin() {
	log.Println("Ordering Gin")
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

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	aws_config := amazonwebservices.StartAws()

	s3_client := amazonwebservices.ConnectS3(aws_config)
	addS3Routes(s3_client, r)

	dynamodb_client := amazonwebservices.ConnectDB(aws_config)
	database.CreateFilesTable(dynamodb_client, "files")
	database.CreateUsersTable(dynamodb_client, "users")
	addDynamoDbRoutes(dynamodb_client, r)

	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	config := fmt.Sprintf(":%s", port)
	log.Printf("Serving Gin at %s:%s", baseUrl, port)
	r.Run(config)
}
