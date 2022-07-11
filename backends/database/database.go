package database

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"

	"fmt"
)

const (
	// RecipeTable is the table name for recipes
	RecipeTable = "recipe"
)

type Client struct {
	dbService dynamodbiface.DynamoDBAPI
}

func New() *Client {
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("AWS_ENDPOINT")

	if region == "" {
		region = "us-west-2"
	}
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	}))
	dbService := dynamodb.New(sess)

	return &Client{
		dbService: dbService,
	}
}

func (client *Client) EnsureTables() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(RecipeTable),
	}

	_, err := client.dbService.CreateTable(input)
	if err != nil {
		log.Default().Printf("error creating table: %v", err)
	}
}

// - MARK: Recipe methods

// SaveRecipe saves a recipe to the DynamoDB Recipe table
func (client *Client) SaveRecipe(recipe Recipe) (*Recipe, error) {
	// generate UUID
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating UUID: %w", err)
	}
	recipe.ID = id.String()

	// marshal recipe
	av, err := dynamodbattribute.MarshalMap(recipe)
	if err != nil {
		return nil, fmt.Errorf("error marshalling recipe item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:         av,
		TableName:    aws.String(RecipeTable),
		ReturnValues: aws.String("ALL_OLD"),
	}

	// save recipe
	_, err = client.dbService.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("error saving recipe: %w", err)
	}

	return &recipe, nil
}

// UpdateRecipe updates an existing recipe
func (client *Client) UpdateRecipe(recipe Recipe, recipeID string) error {
	recipe.ID = recipeID
	av, err := dynamodbattribute.MarshalMap(recipe)
	if err != nil {
		return fmt.Errorf("error marshalling recipe item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:         av,
		TableName:    aws.String(RecipeTable),
		ReturnValues: aws.String("ALL_OLD"),
	}

	_, err = client.dbService.PutItem(input)
	if err != nil {
		return fmt.Errorf("error updating recipe: %w", err)
	}

	return nil
}

// ListAllRecipes returns a list of all recipes as a slice of recipe structs
func (client *Client) ListAllRecipes() ([]Recipe, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(RecipeTable),
	}

	result, err := client.dbService.Scan(params)
	if err != nil {
		return nil, err
	}

	recipes := []Recipe{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &recipes)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

// GetRecipe fetches a recipe by it's ID
func (client *Client) GetRecipe(id string) (*Recipe, error) {
	var recipe *Recipe

	result, err := client.dbService.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(RecipeTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("Could not find recipe with id: " + id)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

// DeleteRecipe deletes a recipe given it's ID
func (client *Client) DeleteRecipe(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(RecipeTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := client.dbService.DeleteItem(input)
	return err
}
