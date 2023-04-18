package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PullRequestInc/go-gpt3"
)

func main() {

	var apiKey string
	var test_name string

	fmt.Println("Please enter your API KEY: ")
	fmt.Scanln(&apiKey)

	if apiKey == "" {
		log.Fatalln("Missing API Key")
	}

	fmt.Println("What is the test you are preparing to? ")
	fmt.Scanln(&test_name)

	var input_text = "Tell me how to prepare for " + test_name + " exam"

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
