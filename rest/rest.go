package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slichlyter12/thyme-apiserver/backends/database"
)

type Client struct {
	Router   *mux.Router
	dbClient *database.Client
}

func New() *Client {
	router := mux.NewRouter()
	router.Use(alwaysJSON)
	router.Use(logging)

	client := &Client{
		Router:   router,
		dbClient: database.New(),
	}

	client.setupRoutes()
	client.dbClient.EnsureTables()
	return client
}

func (client *Client) setupRoutes() {
	apiRouter := client.Router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/status", handleStatus)
	apiRouter.HandleFunc("/recipe", client.handleRecipe)
	apiRouter.HandleFunc("/recipe/{id}", client.handleRecipe)
}

// handles the /status route
func handleStatus(w http.ResponseWriter, r *http.Request) {
	ok := map[string]string{
		"message": "ok",
	}
	okBytes, _ := json.Marshal(ok)
	w.Write(okBytes)
}

// handles the /recipe route
func (client *Client) handleRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch r.Method {
	case "GET":
		if len(vars) > 0 {
			client.getRecipe(w, r, vars["id"])
		} else {
			client.listRecipes(w, r)
		}
		return
	case "POST":
		client.saveRecipe(w, r)
		return
	case "PUT":
		client.updateRecipe(w, r, vars["id"])
		return
	case "DELETE":
		client.deleteRecipe(w, r, vars["id"])
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

// - MARK: Recipe methods

// save a recipe from the body of the method in JSON format
func (client *Client) saveRecipe(w http.ResponseWriter, r *http.Request) {
	// decode request
	decoder := json.NewDecoder(r.Body)
	var recipe database.Recipe
	err := decoder.Decode(&recipe)
	if err != nil {
		writeError(w, "error parsing JSON request", http.StatusBadRequest)
		return
	}

	// save recipe
	savedRecipe, err := client.dbClient.SaveRecipe(recipe)
	if err != nil {
		writeError(
			w,
			"could not save recipe, please try again later",
			http.StatusInternalServerError)
		return
	}

	// encode response
	recipeJSON, err := json.Marshal(savedRecipe)
	if err != nil {
		writeError(w, "could not encode recipe", http.StatusInternalServerError)
		return
	}

	// write response
	writeBytesStatus(w, recipeJSON, http.StatusCreated)
}

// update an existing recipe
func (client *Client) updateRecipe(w http.ResponseWriter, r *http.Request, recipeID string) {
	// get already existing recipe
	oldRecipe, err := client.dbClient.GetRecipe(recipeID)
	if err != nil {
		writeError(w, "could not find recipe with that id", http.StatusNotFound)
		return
	}

	// get updated receipe details
	decoder := json.NewDecoder(r.Body)
	var updatedRecipe database.Recipe
	err = decoder.Decode(&updatedRecipe)
	if err != nil {
		writeError(w, "error parsing json request", http.StatusBadRequest)
		return
	}

	// update recipe
	err = client.dbClient.UpdateRecipe(updatedRecipe, oldRecipe.ID)
	if err != nil {
		writeError(w, "could not update recipe", http.StatusInternalServerError)
		return
	}

	// send response
	w.WriteHeader(http.StatusNoContent)
}

// return a list of all recipes
func (client *Client) listRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := client.dbClient.ListAllRecipes()
	if err != nil {
		writeError(w, "error listing recipes", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(recipes)
	if err != nil {
		writeError(w, "could not marshal recipes", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

// return a singular list with the recipe of the given ID
func (client *Client) getRecipe(w http.ResponseWriter, r *http.Request, id string) {
	recipe, err := client.dbClient.GetRecipe(id)
	if err != nil {
		writeError(w, "could not get recipe with that id", http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(recipe)
	if err != nil {
		writeError(w, "could not marshal recipes", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (client *Client) deleteRecipe(w http.ResponseWriter, r *http.Request, id string) {
	_, err := client.dbClient.GetRecipe(id)
	if err != nil {
		writeError(w, "could not find recipe with that id", http.StatusNotFound)
		return
	}

	err = client.dbClient.DeleteRecipe(id)
	if err != nil {
		writeError(w, "could not delete recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeBytesStatus(w, nil, http.StatusNoContent)
}

// - MARK: Helper Functions

func writeError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error": "` + errorMessage + `"}`))
}

func writeBytesStatus(w http.ResponseWriter, bytes []byte, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write(bytes)
}
