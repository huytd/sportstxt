package main

import (
	"log"
	"net/http"
	"os"

	"huy.rocks/sports/sports"
)

func main() {
	port := "9090"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	handler := sports.NewHandler()

	log.Printf("Starting sportstxt MLB tracker on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
