package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// PexelsResponse represents the response from Pexels API
type PexelsResponse struct {
	Videos []PexelsVideo `json:"videos"`
	Page   int           `json:"page"`
	Total  int           `json:"total"`
}

// PexelsVideo represents a single video from Pexels
type PexelsVideo struct {
	ID           int         `json:"id"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	Duration     int         `json:"duration"`
	Url          string      `json:"url"`
	Image        string      `json:"image"`
	Photographer string      `json:"photographer"`
	VideoFiles   []VideoFile `json:"video_files"`
}

// VideoFile represents video file options from Pexels
type VideoFile struct {
	ID     int    `json:"id"`
	Quality string `json:"quality"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}

// Video represents our internal video structure
type Video struct {
	Title        string `json:"title"`
	Photographer string `json:"photographer"`
	Duration     int    `json:"duration"`
	VideoURL     string `json:"video_url"`
	Thumbnail    string `json:"thumbnail"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

// Pexels API configuration
const (
	PexelsAPIURL = "https://api.pexels.com/videos/search"
)

// Search terms for variety
var searchTerms = []string{
	"nature", "city", "animals", "technology", "people",
	"weather", "water", "mountains", "space", "abstract",
	"food", "travel", "sports", "music", "sunset",
	"ocean", "forest", "sky", "urban", "garden",
	"birds", "flowers", "landscape", "sunrise", "traffic",
	"beach", "snow", "rain", "wind", "clouds",
}

// GetRandomVideo fetches a random video from Pexels API
func GetRandomVideo() (*Video, error) {
	// Pick a random search term
	searchTerm := searchTerms[rand.Intn(len(searchTerms))]

	// Create request to Pexels API (using public endpoint, no auth needed)
	client := &http.Client{Timeout: 10 * time.Second}
	
	url := fmt.Sprintf("%s?query=%s&per_page=80&page=%d", 
		PexelsAPIURL, 
		searchTerm, 
		rand.Intn(5)+1)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching from Pexels: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Pexels API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var pexelsResp PexelsResponse
	err = json.NewDecoder(resp.Body).Decode(&pexelsResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing Pexels response: %w", err)
	}

	// Check if we got any videos
	if len(pexelsResp.Videos) == 0 {
		return nil, fmt.Errorf("no videos found for search term: %s", searchTerm)
	}

	// Pick a random video from results
	randomVideo := pexelsResp.Videos[rand.Intn(len(pexelsResp.Videos))]

	// Find a suitable video file (prefer HD)
	var videoLink string
	for _, vf := range randomVideo.VideoFiles {
		if vf.Quality == "hd" || vf.Quality == "sd" {
			videoLink = vf.Link
			break
		}
	}

	if videoLink == "" && len(randomVideo.VideoFiles) > 0 {
		videoLink = randomVideo.VideoFiles[0].Link
	}

	if videoLink == "" {
		return nil, fmt.Errorf("no valid video file found")
	}

	// Convert to our Video structure
	video := &Video{
		Title:        fmt.Sprintf("%s Video", searchTerm),
		Photographer: randomVideo.Photographer,
		Duration:     randomVideo.Duration,
		VideoURL:     videoLink,
		Thumbnail:    randomVideo.Image,
		Width:        randomVideo.Width,
		Height:       randomVideo.Height,
	}

	return video, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
