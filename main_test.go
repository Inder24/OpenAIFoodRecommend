package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// create mock server and test API connection
func TestOpenAIAPI(t *testing.T) {
	expectedResponse := "Once upon a time, there was a beautiful princess."

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := CompletionResponse{
			Choices: []struct {
				Text string `json:"text"`
			}{
				{
					Text: expectedResponse,
				},
			},
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			t.Errorf("JSON marshaling failed: %v", err)
		}

		w.Write(jsonData)
	}))
	defer server.Close()

	apiURL := server.URL

	// Make the API request
	var receivedResponse string
	completionReq := CompletionRequest{
		Model:       "text-davinci-003",
		Prompt:      prompt,
		MaxTokens:   maxTokens,
		Temperature: 0,
	}

	jsonData, err := json.Marshal(completionReq)
	if err != nil {
		t.Errorf("JSON marshaling failed: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Failed to create HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	var completionResp CompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		t.Errorf("JSON unmarshaling failed: %v", err)
	}

	if len(completionResp.Choices) > 0 {
		receivedResponse = completionResp.Choices[0].Text
	}

	// Compare the expected and received responses
	if receivedResponse != expectedResponse {
		t.Errorf("Unexpected response. Expected: %s, Received: %s", expectedResponse, receivedResponse)
	}
}
