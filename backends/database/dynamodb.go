package database

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	uuid "github.com/nu7hatch/gouuid"

	"fmt"
	"os"
)

var svc *dynamodb.DynamoDB

// Recipe that users can create
type Recipe struct {
	ID          string
	Name        string
	Author      string
	Description string
	Cuisine     string
	ImageName   string
	Ingredients map[string]string
	Steps       []string
}

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

// - MARK: Table Methods

// CreateRecipeTable creates a new table to store recipes in
func CreateRecipeTable() error {
	tableName := "Recipe"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Author"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
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
		return err
	}

	return nil
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

// - MARK: Recipe methods

// SaveRecipe saves a recipe to the DynamoDB Recipe table
func SaveRecipe(recipe Recipe) (string, error) {
	// generate UUID
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Error generating UUID, %s", err.Error())
	}
	recipe.ID = id.String()

	// marshal recipe
	av, err := dynamodbattribute.MarshalMap(recipe)
	if err != nil {
		fmt.Printf("Error marshalling recipe item: %s\n", err.Error())
		return "", err
	}

	tableName := "Recipe"

	input := &dynamodb.PutItemInput{
		Item:         av,
		TableName:    aws.String(tableName),
		ReturnValues: aws.String("ALL_OLD"),
	}

	// save recipe
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Printf("Error saving recipe: %s\n", err.Error())
		return "", err
	}

	return id.String(), nil
}

// ListAllRecipes returns a list of all recipes as a slice of recipe structs
func ListAllRecipes() ([]Recipe, error) {
	tableName := "Recipe"
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		return nil, err
	}

	recipes := []Recipe{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &recipes)
	if err != nil {
		return recipes, err
	}

	return recipes, nil
}

// GetRecipe fetches a recipe by it's ID
func GetRecipe(id string) (Recipe, error) {
	tableName := "Recipe"

	var recipe Recipe

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return recipe, err
	}

	if result.Item == nil {
		return recipe, errors.New("Could not find recipe with id: " + id)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &recipe)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}
