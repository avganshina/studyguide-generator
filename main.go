package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>GPT-3 Demo</title>
		</head>
		<body>
			<form action="/generate" method="POST">
				<label for="api-key">API Key:</label>
				<input type="text" name="api-key" id="api-key"><br><br>
				<label for="test-name">Test Name:</label>
				<input type="text" name="test-name" id="test-name"><br><br>
				<input type="submit" value="Generate">
			</form>
		</body>
		</html>
	`)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.FormValue("api-key")
	testName := r.FormValue("test-name")

	if apiKey == "" {
		http.Error(w, "API Key is required", http.StatusBadRequest)
		return
	}

	if testName == "" {
		http.Error(w, "Test Name is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	inputText := "Tell me how to prepare for " + testName + " exam"

	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			inputText,
		},
		MaxTokens:   gpt3.IntPtr(300),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Fprint(w, resp.Choices[0].Text)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
