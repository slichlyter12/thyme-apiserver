package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/slichlyter12/thyme-apiserver/rest"
)

func main() {
	restClient := rest.New()

	fmt.Println("Serving on 8080...")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(restClient.Router)))
}
