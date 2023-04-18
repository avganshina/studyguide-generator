package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
)

func main() {
	http.HandleFunc("/getresponse", func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.FormValue("apikey")
		testName := r.FormValue("testname")

		if apiKey == "" {
			http.Error(w, "Missing API Key", http.StatusBadRequest)
			return
		}

		if testName == "" {
			http.Error(w, "Missing Test Name", http.StatusBadRequest)
			return
		}

		var inputText = "Tell me how to prepare for " + testName + " exam"

		ctx := context.Background()
		client := gpt3.NewClient(apiKey)

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
			log.Fatalln(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
