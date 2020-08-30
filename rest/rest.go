package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slichlyter12/thyme-apiserver/backends/database"
)

// Router handles all the URLS for the API server
var Router *mux.Router

func init() {
	Router = mux.NewRouter()
	Router.HandleFunc("/init", handleInit)
	Router.HandleFunc("/table", handleTable)
	Router.HandleFunc("/recipe", handleRecipe)
}

// handles the /init route
func handleInit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		setup(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// handles the /table route
func handleTable(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listTables(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// handles the /recipe route
func handleRecipe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listRecipes(w, r)
		return
	case "POST":
		saveRecipe(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// - MARK: Table methods

// creates the Recipe table
func setup(w http.ResponseWriter, r *http.Request) {
	err := database.CreateRecipeTable()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(`OK`))
}

// list all tables in DynamoDB, including those not used
// by this API server
func listTables(w http.ResponseWriter, r *http.Request) {
	tables, derr := database.GetTables()
	if tables == nil {
		w.Write([]byte(`[]`))
		return
	}
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(derr.Error()))
	}
	bytes, err := json.Marshal(tables)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// - MARK: Recipe methods

// save a recipe from the body of the method in JSON format
func saveRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var recipe database.Recipe
	err := decoder.Decode(&recipe)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "error parsing json request"}`))
		return
	}

	id, err := database.SaveRecipe(recipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error saving recipe"}`))
		return
	}

	fmt.Fprintf(w, id)
}

// return a list of all recipes
func listRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recipes, err := database.ListAllRecipes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error listing recipes"}`))
		return
	}

	bytes, err := json.Marshal(recipes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error marshalling recipes"}`))
		return
	}

	w.Write(bytes)
}
