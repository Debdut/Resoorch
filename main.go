package main

import (
	// "embed"

	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/debdut/Resoorch/lib/api/gpt"
)

// // go:embed all:app
// var nextFS embed.FS

// func exaSearchExample(e *exa.Exa) {
// 	query := "hottest AI startups"
// 	response, err := e.Search(query, 2)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("Search results for query '%s':\n", query)
// 	fmt.Printf("Autoprompt: %s\n", response.AutopromptString)
// 	fmt.Printf("Resolved Search Type: %s\n", response.ResolvedSearchType)
// 	fmt.Printf("Request ID: %s\n\n", response.RequestID)

// 	for i, result := range response.Results {
// 		fmt.Printf("Result %d:\n", i+1)
// 		fmt.Printf("  Title: %s\n", result.Title)
// 		fmt.Printf("  URL: %s\n", result.URL)
// 		fmt.Printf("  Score: %.4f\n", result.Score)
// 		fmt.Printf("  Published Date: %s\n", result.PublishedDate)
// 		fmt.Printf("  Author: %s\n", result.Author)
// 		fmt.Printf("  Text: %s\n\n", result.Text)
// 	}
// }

func gptQueryExample(g *gpt.GPT) {
	responseFormat := gpt.ResponseFormat{
		Type: "json_schema",
		JSONSchema: gpt.JSONSchema{
			Name:   "math_response",
			Strict: true,
			Schema: gpt.Schema{
				Type: "object",
				Properties: map[string]gpt.Property{
					"steps": {
						Type: "array",
						Items: map[string]interface{}{
							"type": "object",
							"properties": map[string]gpt.Property{
								"explanation": {Type: "string"},
								"output":      {Type: "string"},
							},
							"required":             []string{"explanation", "output"},
							"additionalProperties": false,
						},
					},
					"final_answer": {Type: "string"},
				},
				Required:             []string{"steps", "final_answer"},
				AdditionalProperties: false,
			},
		},
	}
	messages := []gpt.Message{
		{Role: "system", Content: "You are a helpful math tutor."},
		{Role: "user", Content: "solve 8x + 31 = 2"},
	}
	response, err := g.Call(&messages, &responseFormat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		if response != nil {
			fmt.Printf(response.Error.Message)
		}
		return
	}
	fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// e, err := exa.InitExa(0)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }
	// exaSearchExample(e)

	g, err := gpt.InitGPT(0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	gptQueryExample(g)
}
