package ai

import (
	"context"
	"fmt"
	config "kesbekes/Config"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type AI struct {
	AIModel *genai.GenerativeModel
}

func NewAI() *AI {
	model := config.NewAI()
	return &AI{
		AIModel: model,
	}
}

// printResponse prints the response content parts for debugging purposes
func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}

// IsPreferred checks if the txt is preferred based on the AI model and preferences
func (ai *AI) IsPreferred(txt string, preferences []string) (bool, error) {
	fmt.Println("this is the model ", ai.AIModel)

	preferencesList := strings.Join(preferences, ", ")
	query := "Does the text '" + txt + "' match any of these preferences: " + preferencesList

	ctx := context.Background()
	resp, err := ai.AIModel.GenerateContent(ctx, genai.Text(query))
	if err != nil {
		return false, err
	}

	// Print the response for debugging
	printResponse(resp)

	//// Check each candidate's content for keywords indicating preference
	// for _, cand := range resp.Candidates {
	// 	if cand.Content != nil {
	// 		for _, part := range cand.Content.Parts {
	// 			lowerPart := strings.ToLower(part)
	// 			if strings.Contains(lowerPart, "yes") || strings.Contains(lowerPart, "matches") {
	// 				return true, nil
	// 			}
	// 		}
	// 	}
	// }

	// If no match was found
	return false, nil
}
