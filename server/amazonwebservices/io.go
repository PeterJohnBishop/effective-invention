package amazonwebservices

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func UploadFile(client *s3.Client, filename string, fileContent multipart.File) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   fileContent, // Directly passing io.Reader
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, filename)
	return fileURL, nil
}

func DownloadFile(client *s3.Client, filename string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	fileKey := filename

	expiration := time.Duration(5) * time.Minute

	presignClient := s3.NewPresignClient(client)
	presignedURL, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return presignedURL.URL, nil
}
