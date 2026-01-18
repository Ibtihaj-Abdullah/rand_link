package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleRandomLink(w http.ResponseWriter, r *http.Request) {
	// Get random video from Pexels API
	video, err := GetRandomVideo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching video: %v", err), http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}
