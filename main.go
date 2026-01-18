package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("static")))

	// API endpoint for random videos
	http.HandleFunc("/api/random", handleRandomLink)

	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
