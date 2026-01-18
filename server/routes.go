package server

import (
	"effective-invention/server/amazonwebservices"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func addS3Routes(s3client *s3.Client, dynamodbClient *dynamodb.Client, r *gin.Engine) {
	r.POST("/upload", amazonwebservices.HandleFileUpload(s3client))
	r.GET("/download/link/:filename", amazonwebservices.HandleFileDOwnloadLink(s3client))
	r.GET("/download/:filename", amazonwebservices.HandleFileDownloadStream(s3client))
	r.POST("/user/upload", amazonwebservices.HandleUploadUserFile(dynamodbClient, s3client))
	r.GET("/download/qrlink/:filename", amazonwebservices.HandleFileDOwnloadLinkQR(s3client))
}

func addDynamoDbRoutes(s3client *s3.Client, dynamodbClient *dynamodb.Client, r *gin.Engine) {
	r.POST("/users/new", amazonwebservices.HandleUserCreation(dynamodbClient))
	r.POST("/users/login", amazonwebservices.HandleAuthentication(dynamodbClient))
	r.GET("/users/all", amazonwebservices.HandleGetAllUsers(dynamodbClient))
	r.GET("/users/id/:id", amazonwebservices.HandleGetUserById(dynamodbClient))
	r.PUT("/users/update", amazonwebservices.HandleUpdateUser(dynamodbClient))
	r.PUT("/users/update/password", amazonwebservices.HandleUpdateUserPassword(dynamodbClient))
	r.DELETE("/users/id/:id", amazonwebservices.HandleDeleteUserById(dynamodbClient))

	r.GET("/users/files", amazonwebservices.HandleGetUserFiles(dynamodbClient))
	r.DELETE("/users/files/:id", amazonwebservices.HandleDeleteUserFileById(dynamodbClient, s3client))
}

func addRekognitionRoutes(client *rekognition.Client, r *gin.Engine) {
	r.GET("/analysis/:file", amazonwebservices.HandleFacialAnalysis(client))
}
