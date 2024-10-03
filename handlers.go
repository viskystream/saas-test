package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
				AppData:    getUserData(stream.Token),
				ViewTokens: make(map[string]ViewTokenResponse),
			}

			for _, viewToken := range stream.ViewTokens {
				viewTokenResponse := ViewTokenResponse{
					Stop:    false,
					AppData: getUserData(viewToken),
				}
				streamResponse.ViewTokens[viewToken.Value] = viewTokenResponse
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
		"user.name":  "Unknown User",
		"user.scope": "viewer", // Default to viewer scope for unknown tokens
	}
}

func getPrivateKey(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/api/ls/v1/key/%s?token=%s", os.Getenv("BACKEND_ENDPOINT"), r.URL.Query().Get("user"), os.Getenv("TOKEN"))

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func getLiveStreams(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/api/ls/v1/live?token=%s", os.Getenv("BACKEND_ENDPOINT"), os.Getenv("TOKEN"))

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func getAuthToken(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/auth/v1/access-tokens", os.Getenv("BACKEND_ENDPOINT"))
	url = strings.Replace(url, "umbrella.", "", 1)

	var requestBody interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var responseData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(responseData)
}

func getViewersWatching(w http.ResponseWriter, r *http.Request) {
	callId := r.URL.Query().Get("callId")

	if callId == "" {
		http.Error(w, "A valid callId is required as a query parameter.", http.StatusBadRequest)
		return
	}

	viewers, exists := viewersByCallId[callId]
	if !exists {
		viewers = []string{}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"callId": callId, "viewers": viewers})
}

func handleViewerJoined(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CallId string `json:"callId"`
		PeerId string `json:"peerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, exists := viewersByCallId[request.CallId]; !exists {
		viewersByCallId[request.CallId] = []string{}
	}

	viewersByCallId[request.CallId] = append(viewersByCallId[request.CallId], request.PeerId)

	message := fmt.Sprintf("Viewer %s joined call %s", request.PeerId, request.CallId)
	streamHub.broadcast <- []byte(message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Viewer joined successfully"})
}

func handleViewerLeft(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CallId string `json:"callId"`
		PeerId string `json:"peerId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if viewers, exists := viewersByCallId[request.CallId]; exists {
		updatedViewers := []string{}
		for _, viewer := range viewers {
			if viewer != request.PeerId {
				updatedViewers = append(updatedViewers, viewer)
			}
		}
		viewersByCallId[request.CallId] = updatedViewers
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Viewer %s removed from call %s.", request.PeerId, request.CallId)})
	} else {
		http.Error(w, fmt.Sprintf("CallId %s not found.", request.CallId), http.StatusNotFound)
	}
}

func handleBroadcastEnded(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CallId string `json:"callId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if _, exists := viewersByCallId[request.CallId]; exists {
		delete(viewersByCallId, request.CallId)
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Broadcast %s has ended and viewers cleared.", request.CallId)})
	} else {
		http.Error(w, fmt.Sprintf("CallId %s not found or already cleared.", request.CallId), http.StatusNotFound)
	}
}
