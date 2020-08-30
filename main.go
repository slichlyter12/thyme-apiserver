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
	message, derr := database.CreateRecipeTable()
	if derr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(derr.Error()))
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		w.Write(bytes)
	}
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
	w.Write(bytes)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/init", setup).Methods("POST")
	router.HandleFunc("/table", listTables).Methods("GET")

	fmt.Println("Serving on 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
