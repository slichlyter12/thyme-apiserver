package database

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, nil
}

func (m *mockDynamoDBClient) Scan(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{}, nil
}

func (m *mockDynamoDBClient) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{}, nil
}

func newMockClient() *Client {
	return &Client{
		dbService: &mockDynamoDBClient{},
	}
}

func TestSaveRecipe(t *testing.T) {
	mockClient := newMockClient()
	recipe := Recipe{
		Name:   "Snickerdoodle Cookies",
		Author: "Gran",
	}
	_, err := mockClient.SaveRecipe(recipe)
	if err != nil {
		t.Errorf("Error saving recipe: %s", err.Error())
	}
}

func TestUpdateRecipe(t *testing.T) {
	mockClient := newMockClient()
	recipe := Recipe{
		Name:   "Butternut Squash Soup",
		Author: "Aleksa Wood",
	}

	savedRecipe, err := mockClient.SaveRecipe(recipe)
	if err != nil {
		t.Errorf("Error saving recipe: %s", err.Error())
	}

	newRecipe := recipe
	newRecipe.Name = "Better Butternut Sqash Soup"

	err = mockClient.UpdateRecipe(newRecipe, savedRecipe.ID)
	if err != nil {
		t.Errorf("Error updating recipe: %s", err.Error())
	}
}

func TestGetRecipeById(t *testing.T) {
	mockClient := newMockClient()
	recipe := Recipe{
		Name:   "Roasted Carrots",
		Author: "Sam Lichlyter",
	}

	savedRecipe, err := mockClient.SaveRecipe(recipe)
	if err != nil {
		t.Errorf("Error saving recipe: %s", err.Error())
	}

	requestedRecipe, err := mockClient.GetRecipe(savedRecipe.ID)
	if err != nil {
		t.Errorf("Error getting recipe: %s", err.Error())
	}

	recipe.ID = savedRecipe.ID
	if recipe.ID != requestedRecipe.ID || recipe.Name != requestedRecipe.Name || recipe.Author != requestedRecipe.Author {
		t.Errorf("Recipes are not the same:\n%+v\n%+v", recipe, requestedRecipe)
	}
}

// func TestGetRecipeByInvalidId(t *testing.T) {
// 	recipe := Recipe{
// 		Name:   "Steak",
// 		Author: "Sam Lichlyter",
// 	}

// 	_, err := SaveRecipe(recipe)
// 	if err != nil {
// 		t.Errorf("Error saving recipe: %s", err.Error())
// 	}

// 	_, err = GetRecipe("1234")
// 	if err == nil {
// 		t.Error("Found recipe with invalid ID")
// 	}
// }

// func TestListAllRecipes(t *testing.T) {
// 	recipes, err := ListAllRecipes()
// 	if err != nil {
// 		t.Errorf("Error listing all recipes: %s", err.Error())
// 	}

// 	if len(recipes) == 0 {
// 		t.Errorf("There are no recipes in the database")
// 	}
// }
