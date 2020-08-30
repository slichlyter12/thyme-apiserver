package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"fmt"
	"os"
)

var svc *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-2"
	}
	endpoint := os.Getenv("AWS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	if session, err := session.NewSession(&aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	}); err != nil {
		fmt.Print("Error creating AWS session", err.Error())
	} else {
		svc = dynamodb.New(session)
	}
}

// CreateRecipeTable creates a new table to store recipes in
func CreateRecipeTable() (string, error) {
	tableName := "Recipe"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Author"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Author"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		return "Error creating Recipe table", err
	}

	return "created Recipe Table", nil
}

// GetTables returns a list of all tables in the DB
func GetTables() ([]string, error) {
	input := &dynamodb.ListTablesInput{}

	var tableNames []string
	for {
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
					return nil, aerr
				default:
					fmt.Println(aerr.Error())
					return nil, aerr
				}
			} else {
				fmt.Println(err.Error())
				return nil, err
			}
		}

		for _, n := range result.TableNames {
			tableNames = append(tableNames, *n)
		}

		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return tableNames, nil
}

// DeleteTable deletes the specified table
func DeleteTable(tableName string) error {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteTable(params)
	if err != nil {
		return err
	}

	return nil
}
