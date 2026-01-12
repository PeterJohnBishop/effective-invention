package server

import (
	"effective-invention/server/amazonwebservices"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func addS3Routes(client *s3.Client, r *gin.Engine) {
	r.POST("/upload", amazonwebservices.HandleFileUpload(client))
	r.GET("/download/:filename", amazonwebservices.HandleFileDOwnload(client))
}

func addDynamoDbRoutes(client *dynamodb.Client, r *gin.Engine) {
	r.POST("/users/new", amazonwebservices.HandleUserCreation(client))
	r.POST("/users/login", amazonwebservices.HandleAuthentication(client))
	r.GET("/users/all", amazonwebservices.HandleGetAllUsers(client))
	r.GET("/users/id/:id", amazonwebservices.HandleGetUserById(client))
	r.PUT("/users/update", amazonwebservices.HandleUpdateUser(client))
	r.PUT("/users/update/password", amazonwebservices.HandleUpdateUserPassword(client))
	r.DELETE("/users/id/:id", amazonwebservices.HandleDeleteUserById(client))
}
