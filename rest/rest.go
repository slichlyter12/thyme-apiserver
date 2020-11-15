package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slichlyter12/thyme-apiserver/backends/database"
)

// Router handles all the URLS for the API server
var Router *mux.Router

func init() {
	Router = mux.NewRouter()
	apiRouter := Router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/status", handleStatus)
	apiRouter.HandleFunc("/init", handleInit)
	apiRouter.HandleFunc("/table", handleTable)
	apiRouter.HandleFunc("/recipe", handleRecipe)
	apiRouter.HandleFunc("/recipe/{id}", handleRecipe)
}

// handles the /status route
func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(`OK`))
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
	vars := mux.Vars(r)

	switch r.Method {
	case "GET":
		if len(vars) > 0 {
			getRecipe(w, r, vars["id"])
		} else {
			listRecipes(w, r)
		}
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
		return
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
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"id": "` + id + `"}`))
}

// return a list of all recipes
func listRecipes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recipes, err := database.ListAllRecipes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error listing recipes", "error": "` + err.Error() + `"}`))
		return
	}

	bytes, err := json.Marshal(recipes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error marshalling recipes", "error": "` + err.Error() + `"}`))
		return
	}

	w.Write(bytes)
}

// return a singular list with the recipe of the given ID
func getRecipe(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "application/json")

	recipe, err := database.GetRecipe(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	bytes, err := json.Marshal(recipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "error marshalling recipe"}`))
	}

	w.Write(bytes)
}
