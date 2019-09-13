package main

import (
	"context"
	"fmt"

	language "cloud.google.com/go/language/apiv1"
	"google.golang.org/grpc"
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
	client := language.NewClient(ctx, grpc.WithInsecure())
	text := "are you crazy!!! this was supposed to be a kid friendly conference!! "
	analysis, err := analyzeSentiment(ctx, client, text)

	fmt.Println(analysis)
}
