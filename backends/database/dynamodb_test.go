package database

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	err := CreateRecipeTable()
	if err != nil {
		t.Errorf("Error creating table: %s", err.Error())
	}
}

func TestCreateTableError(t *testing.T) {
	err := CreateRecipeTable()
	if err == nil {
		t.Errorf("Failed to produce an error")
	}
}

func TestGetTables(t *testing.T) {
	tables, _ := GetTables()
	expectedNumberOfTables := 1
	if len(tables) != 1 {
		t.Errorf("There are %d tables instead of %d", len(tables), expectedNumberOfTables)
	}
}

func TestSaveRecipe(t *testing.T) {
	recipe := Recipe{
		Name:   "Snickerdoodle Cookies",
		Author: "Gran",
	}
	err := SaveRecipe(recipe)
	if err != nil {
		t.Errorf("Error saving recipe: %s", err.Error())
	}
}

func TestListAllRecipes(t *testing.T) {
	recipes, err := ListAllRecipes()
	if err != nil {
		t.Errorf("Error listing all recipes: %s", err.Error())
	}

	expectedNumberOfRecipes := 1
	if len(recipes) != expectedNumberOfRecipes {
		t.Errorf("There are %d recipes instead of %d", len(recipes), expectedNumberOfRecipes)
	}
}

func TestDeleteTable(t *testing.T) {
	err := DeleteTable("Recipe")
	if err != nil {
		t.Errorf("Error deleting table: %s", err.Error())
	}
}
