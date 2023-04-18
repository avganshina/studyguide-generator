package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
)

type Request struct {
	TestName string `json:"testName"`
	ApiKey   string `json:"apiKey"`
}

type Response struct {
	Answer string `json:"answer"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ApiKey == "" {
		http.Error(w, "Missing API key", http.StatusBadRequest)
		return
	}

	if req.TestName == "" {
		http.Error(w, "Missing test name", http.StatusBadRequest)
		return
	}

	inputText := "Tell me how to prepare for " + req.TestName + " exam"

	ctx := context.Background()
	client := gpt3.NewClient(req.ApiKey)

	var responseText string
	err = client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			inputText,
		},
		MaxTokens:   gpt3.IntPtr(300),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		responseText = resp.Choices[0].Text
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{Answer: responseText}
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
