package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type CompletionRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type CompletionResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func generateResponse(input string, model string, apiKey string) (string, error) {
	// Create a new completion request
	completionRequest := &CompletionRequest{
		Model:     model,
		Prompt:    fmt.Sprintf("%s\nMaster Control Program:", input),
		MaxTokens: 1000,
	}
	payload, err := json.Marshal(completionRequest)
	if err != nil {
		return "", err
	}

	// Send the completion request to the OpenAI API endpoint
	url := "https://api.openai.com/v1/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the completion response
	completionResponse := &CompletionResponse{}
	err = json.NewDecoder(resp.Body).Decode(completionResponse)
	if err != nil {
		return "", err
	}

	// Return the text response
	return completionResponse.Choices[0].Text, nil
}

func main() {
	// Prompt user for API key
	fmt.Print("Enter your OpenAI API key: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	apiKey := scanner.Text()

	// Initialize scanner for user input
	scanner = bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(string("\033[34m"), "User: ")
		scanner.Scan()
		userInput := scanner.Text()
		if strings.ToLower(userInput) == "quit" {
			break
		}
		response, err := generateResponse(userInput, "text-davinci-003", apiKey)
		if err != nil {
			fmt.Println("Error generating response:", err)
		} else {
			fmt.Println(string("\033[31m"), "Master Control Program:", string("\033[31m"), " "+response)

		}
	}
}