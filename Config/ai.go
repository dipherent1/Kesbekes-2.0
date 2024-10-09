package config

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func NewAI() *genai.GenerativeModel {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(GeminAPIKey))
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	return model

}
