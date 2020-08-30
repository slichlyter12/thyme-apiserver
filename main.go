package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/slichlyter12/thyme-apiserver/backends/database"
)

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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/init", setup).Methods("POST")
	router.HandleFunc("/table", listTables).Methods("GET")
	router.HandleFunc("/recipe", saveRecipe).Methods("POST")

	fmt.Println("Serving on 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
