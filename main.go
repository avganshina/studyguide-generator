package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

func main() {
	apiKey := os.Getenv("API_KEY")

	if apiKey == "" {
		log.Fatalln("Missing API Key")
	}

	var test_name string
	fmt.Println("What is the test you are preparing for?")
	fmt.Scanln(&test_name)

	input_text := "Tell me how to prepare for " + test_name + " exam"

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			input_text,
		},
		MaxTokens:   gpt3.IntPtr(300),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})

	if err != nil {
		log.Fatalln(err)
	}
}
