package amazonwebservices

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

func ConnectRekognition(cfg aws.Config) *rekognition.Client {
	rekognitionClient := rekognition.NewFromConfig(cfg)
	if rekognitionClient == nil {
		log.Fatalf("Error connecting to Rekognition.")
	}
	log.Printf("Connected to Rekognition\n")
	return rekognitionClient
}

type FaceAnalysis struct {
	AgeRange      string
	Gender        string
	Emotions      []string
	Smile         bool
	Eyeglasses    bool
	Sunglasses    bool
	Beard         bool
	Mustache      bool
	EyesOpen      bool
	MouthOpen     bool
	Confidence    float32
	LandmarkCount int
}

func AnalyzeFace(client *rekognition.Client, fileKey string) ([]FaceAnalysis, error) {
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	// DEBUG LOGS
	log.Printf("DEBUG: Looking in Bucket: %s", bucketName)
	log.Printf("DEBUG: Looking for Key: %s", fileKey)
	log.Printf("DEBUG: Using Region: %s", region)

	input := &rekognition.DetectFacesInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(fileKey),
			},
		},
		Attributes: []types.Attribute{types.AttributeAll},
	}

	result, err := client.DetectFaces(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("rekognition error: %w", err)
	}

	var analyses []FaceAnalysis

	for _, detail := range result.FaceDetails {
		analysis := FaceAnalysis{
			AgeRange:   fmt.Sprintf("%d-%d", *detail.AgeRange.Low, *detail.AgeRange.High),
			Gender:     string(detail.Gender.Value),
			Confidence: *detail.Confidence,

			Smile:      detail.Smile.Value,
			Eyeglasses: detail.Eyeglasses.Value,
			Sunglasses: detail.Sunglasses.Value,
			Beard:      detail.Beard.Value,
			Mustache:   detail.Mustache.Value,
			EyesOpen:   detail.EyesOpen.Value,
			MouthOpen:  detail.MouthOpen.Value,

			LandmarkCount: len(detail.Landmarks),
		}

		for _, e := range detail.Emotions {
			if *e.Confidence > 50.0 {
				analysis.Emotions = append(analysis.Emotions, string(e.Type))
			}
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

type ComparisonResult struct {
	IsMatch        bool
	Similarity     float32
	FaceLocation   *types.BoundingBox
	UnmatchedCount int
}

func CompareTwoFaces(client *rekognition.Client, sourceKey, targetKey string) (ComparisonResult, error) {
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	input := &rekognition.CompareFacesInput{
		SourceImage: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(sourceKey),
			},
		},
		TargetImage: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucketName),
				Name:   aws.String(targetKey),
			},
		},
		SimilarityThreshold: aws.Float32(80.0), // 80% is usually the standard for "likely same person"
	}

	result, err := client.CompareFaces(context.TODO(), input)
	if err != nil {
		return ComparisonResult{}, fmt.Errorf("comparison failed: %w", err)
	}

	res := ComparisonResult{
		UnmatchedCount: len(result.UnmatchedFaces),
	}

	if len(result.FaceMatches) > 0 {
		topMatch := result.FaceMatches[0]
		res.IsMatch = true
		res.Similarity = *topMatch.Similarity
		res.FaceLocation = topMatch.Face.BoundingBox
	}

	return res, nil
}
