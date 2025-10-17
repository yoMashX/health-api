package main

import (
	"log"
	"net/http"

	"health-api/internal/api"
	"health-api/internal/database"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	router := api.NewRouter(db)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}