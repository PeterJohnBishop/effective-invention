package amazonwebservices

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ConnectS3(cfg aws.Config) *s3.Client {
	s3Client := s3.NewFromConfig(cfg)
	_, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("unable to load S3 buckets, %v", err)
	}
	log.Printf("Connected to S3\n")
	return s3Client
}
