package server

import (
	"effective-invention/server/amazonwebservices"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func addS3Routes(client *s3.Client, r *gin.Engine) {
	r.POST("/upload", amazonwebservices.HandleFileUpload(client))
	r.GET("/download/:filename", amazonwebservices.HandleFileDOwnload(client))
}
