package main

import (
	"context"
	"fmt"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}
func main() {
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	text := "I love Ice cream! I hate mosquitos"
	analysis, err := analyzeSentiment(ctx, client, text)
	sentences := analysis.GetSentences()
	for _, s := range sentences {
		fmt.Printf("sentence: %v \nlanguage: %v \nsentiment: %v\n", s.GetText(), analysis.GetLanguage(), s.GetSentiment())
	}
}
