package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
)

func main() {
	http.HandleFunc("/prepare-test", prepareTestHandler)
	http.ListenAndServe(":8080", nil)
}

func prepareTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// Render the form template
		tmpl, err := template.ParseFiles("templates/prepare-test.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle the form submission
	apiKey := r.FormValue("apikey")
	examName := r.FormValue("examname")

	if apiKey == "" || examName == "" {
		http.Error(w, "Missing API key or exam name", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	var text strings.Builder
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			"Tell me how to prepare for " + examName + " exam",
		},
		MaxTokens:   gpt3.IntPtr(300),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Fprint(&text, resp.Choices[0].Text)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save text to disk
	textStr := text.String()
	err = os.WriteFile("exam-preparation.txt", []byte(textStr), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve the text file as a download
	w.Header().Set("Content-Disposition", "attachment; filename=exam-preparation.txt")
	w.Header().Set("Content-Type", "text/plain")
	http.ServeFile(w, r, "exam-preparation.txt")

	// Delete the text file from disk
	err = os.Remove("exam-preparation.txt")
	if err != nil {
		fmt.Println("Error deleting text file:", err)
	}
}
