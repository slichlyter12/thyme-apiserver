package database

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	_, err := CreateRecipeTable()
	if err != nil {
		t.Errorf("Error creating table: %s", err.Error())
	}
}

func TestCreateTableError(t *testing.T) {
	_, err := CreateRecipeTable()
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

func TestDeleteTable(t *testing.T) {
	err := DeleteTable("Recipe")
	if err != nil {
		t.Errorf("Error deleting table: %s", err.Error())
	}
}
