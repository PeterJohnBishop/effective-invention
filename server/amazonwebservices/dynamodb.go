package amazonwebservices

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func ConnectDB(cfg aws.Config) *dynamodb.Client {
	dynamoClient := dynamodb.NewFromConfig(cfg)
	_, err := GetTables(dynamoClient)
	if err != nil {
		log.Fatalf("Error connecting to DynamoDB.")
	}
	log.Printf("Connected to DynamoDB\n")
	return dynamoClient
}

func GetTables(client *dynamodb.Client) ([]string, error) {

	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}
	return result.TableNames, nil
}
