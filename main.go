package main

import (
	"example/FirstApi/handlers"
	"example/FirstApi/models"
	"log"
	"net/http"
	"os"
)

func main() {
	handlers.Items = append(handlers.Items, models.Item{ID: 1, Name: "Item 1", Price: 100})
	handlers.Items = append(handlers.Items, models.Item{ID: 2, Name: "Item 2", Price: 200})

	http.HandleFunc("/items", handlers.GetItems)
	http.HandleFunc("/items/", handlers.ItemHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	log.Printf("listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
