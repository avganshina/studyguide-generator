package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
)

type Request struct {
	APIKey   string `json:"api_key"`
	TestName string `json:"test_name"`
}

type Response struct {
	Result string `json:"result"`
}

func Gpt3Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.APIKey == "" || req.TestName == "" {
		http.Error(w, "Missing API key or test name", http.StatusBadRequest)
		return
	}

	inputText := "Tell me how to prepare for " + req.TestName + " exam"

	ctx := context.Background()
	client := gpt3.NewClient(req.APIKey)

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

	res := Response{Result: responseText}

	resJSON, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resJSON)
}
