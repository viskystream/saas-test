package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var viewersByCallId = make(map[string][]string)

var mockUsers = map[string]map[string]string{
	"broadcaster_token_123": {
		"user.id":    "broadcaster123",
		"user.name":  "John Broadcaster",
		"user.scope": "broadcaster",
	},
	"viewer_token_456": {
		"user.id":    "viewer456",
		"user.name":  "Jane Viewer",
		"user.scope": "viewer",
	},
	"viewer_token_789": {
		"user.id":    "viewer789",
		"user.name":  "Bob Watcher",
		"user.scope": "viewer",
	},
}

type WebhookRequest struct {
	Programs map[string]Program `json:"programs"`
}

type Program struct {
	Streams map[string]Stream `json:"streams"`
}

type Stream struct {
	Token      Token   `json:"token"`
	ViewTokens []Token `json:"viewTokens"`
}

type Token struct {
	Value  string `json:"value"`
	Type   string `json:"type"`
	Action string `json:"action"`
}

// WebhookResponse represents the response structure
type WebhookResponse struct {
	Programs map[string]ProgramResponse `json:"programs"`
}

type ProgramResponse struct {
	Stop     bool                      `json:"stop"`
	NeedAuth bool                      `json:"needAuth"`
	Streams  map[string]StreamResponse `json:"streams"`
}

type StreamResponse struct {
	Stop       bool                         `json:"stop"`
	NeedAuth   bool                         `json:"needAuth"`
	Token      string                       `json:"token"`
	AppData    map[string]string            `json:"appData"`
	ViewTokens map[string]ViewTokenResponse `json:"viewTokens"`
}

type ViewTokenResponse struct {
	Stop    bool              `json:"stop"`
	AppData map[string]string `json:"appData"`
}

// Add this new function to handle webhook requests
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Print the raw payload
	log.Printf("Received webhook payload: %s", string(body))

	// Restore the request body for further processing
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var request WebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := validateAndProcessWebhook(request)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateAndProcessWebhook(request WebhookRequest) WebhookResponse {
	response := WebhookResponse{
		Programs: make(map[string]ProgramResponse),
	}

	for programID, program := range request.Programs {
		programResponse := ProgramResponse{
			Stop:     false,
			NeedAuth: true,
			Streams:  make(map[string]StreamResponse),
		}

		for streamID, stream := range program.Streams {
			streamResponse := StreamResponse{
				Stop:       false,
				NeedAuth:   true,
				Token:      validateToken(stream.Token),
				AppData:    make(map[string]string),
				ViewTokens: make(map[string]ViewTokenResponse),
			}

			switch stream.Token.Action {
			case "joining":
			case "creating":
				streamResponse.AppData = getUserData(stream.Token) // Value is the username
			case "polling":
				for _, viewToken := range stream.ViewTokens {
					viewTokenResponse := ViewTokenResponse{
						Stop:    false,
						AppData: getUserData(viewToken), // Value is the username
					}
					streamResponse.ViewTokens[viewToken.Value] = viewTokenResponse
				}
			}

			programResponse.Streams[streamID] = streamResponse
		}

		response.Programs[programID] = programResponse
	}

	return response
}

func validateToken(token Token) string {
	// echo the token
	return token.Value
}

func getUserData(token Token) map[string]string {
	// Check if the token value exists in our mock users
	if userData, exists := mockUsers[token.Value]; exists {
		return userData
	}

	// If the token is not found, return a default set of data
	return map[string]string{
		"user.id":    "unknown",
		"user.name":  token.Value,
		"user.scope": "viewer", // Default to viewer scope for unknown tokens
	}
}
