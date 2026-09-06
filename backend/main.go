package main

import (
	"fmt"
	"log"
	"net/http"
	"nimbus/api"
)

func main() {
	r := api.SetupRouter()

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
