package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slichlyter12/thyme-apiserver/backends/database"
)

// Router is the router used by the rest package
var Router *mux.Router

func init() {
	Router = mux.NewRouter()
	Router.HandleFunc("/init", handleInit)
	Router.HandleFunc("/table", handleTable)
	Router.HandleFunc("/recipe", handleRecipe)
}

func handleInit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		setup(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handleTable(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listTables(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

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

func setup(w http.ResponseWriter, r *http.Request) {
	err := database.CreateRecipeTable()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(`OK`))
}

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

	err = database.SaveRecipe(recipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error saving recipe"}`))
		return
	}

	w.Write([]byte(`{"message": "saved recipe"}`))
}

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
