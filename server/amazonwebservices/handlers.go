package amazonwebservices

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func HandleFileUpload(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file") // "file" is the key in the form-data
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
			return
		}
		defer file.Close()

		url, err := UploadFile(client, header.Filename, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded successfully",
			"url":     url,
		})
	}
}

func HandleFileDOwnload(client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		if filename == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
			return
		}

		url, err := DownloadFile(client, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File download link generated successfully",
			"expires": "5 minutes",
			"url":     url,
		})
	}
}
