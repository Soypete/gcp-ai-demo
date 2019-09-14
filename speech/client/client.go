package main

import (
	"context"
	"fmt"
	"os"

	proto "github.com/Soypete/gcp-ai-demo/proto"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	option := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:3030", option)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewTranscripterClient(conn)

	p, err := client.SpeechToText(ctx, &proto.Voice{VoiceLocation: os.Getenv("LOCATION"), Bucket: "phone-call-analysis", TextLocation: "voice-mails", ProjectID: os.Getenv("PROJECT_ID")})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Text: %s\nConfidence: %f.02\n", p.GetText(), p.GetConfidence())
}
