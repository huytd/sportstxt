package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"huy.rocks/sports/sports"
)

func main() {
	_ = godotenv.Load()

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
