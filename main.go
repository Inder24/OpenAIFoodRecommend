package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiURL = "https://api.openai.com/v1/completions"
	apiKey = "<ADD-TOKEN>"
	prompt = `Hello, my name is John. I prefer sweet and umami flavors, and I'm currently located in the United States. It's evening here, and I'm feeling relaxed. I'm looking for a dinner recipe for four people. I enjoy Italian and Mexican cuisines but need the recipe to be gluten-free. I'd prefer not to have broccoli in it. By the way, I don't have any religious restrictions on my diet. 

the output should be in the following format only :

{
    "FoodItem": "Chicken Enchiladas",
    "RecipeLink": "www.example-recipe-link.com",
    "PreparationTime": "60 minutes",
    "Difficulty": "Intermediate",
    "Servings": "2",
    "Ingredients": [
        "Chicken",
        "Tortillas",
        "Cheese",
        "Enchilada Sauce",
        "Onions",
        "Bell Peppers"
    ],
    "ImageLink": "www.example-google-image-link.com",
    "NutritionalInformation": {
        "Calories": "550 kcal",
        "Protein": "30 g",
        "Fat": "25 g",
        "Carbohydrates": "50 g"
    }
}`
	maxTokens = 3500
)

type CompletionRequest struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

type RecipeResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int            `json:"created"`
	Model   string         `json:"model"`
	Choices []RecipeChoice `json:"choices"`
	Usage   UsageInfo      `json:"usage"`
}

type RecipeChoice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	completionReq := CompletionRequest{
		Model:       "text-davinci-003",
		Prompt:      prompt,
		MaxTokens:   maxTokens,
		Temperature: 0,
	}

	jsonData, err := json.Marshal(completionReq)
	if err != nil {
		fmt.Println("JSON marshaling failed:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var completionResp RecipeResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		fmt.Println("JSON unmarshaling failed:", err)
		return
	}

	if len(completionResp.Choices) > 0 {
		fmt.Printf("Response: %+v", completionResp)
	} else {
		fmt.Println("No response received")
	}
}
