package main

import (
	"context"
	"fmt"
	"log"
	"net"

	proto "github.com/Soypete/gcp-ai-demo/proto"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SpeechToText(ctx context.Context, in *proto.Voice) (*proto.Text, error) {

	script, conf, err := transcribe(ctx, in.GetVoiceLocation())
	if err != nil {
		panic(err)
	}
	return &proto.Text{Text: script, Confidence: conf}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:3030")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterTranscripterServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
func transcribe(ctx context.Context, filename string) (script string, conf float32, err error) {
	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_ENCODING_UNSPECIFIED,
			SampleRateHertz: 8000,
			LanguageCode:    "en-US",
			Model:           "phone_call",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: filename},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}
	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			script = fmt.Sprintf("%s%s\n", script, alt.Transcript)
			conf = conf + alt.Confidence
		}
	}
	conf = conf / float32(len(resp.Results))
	return script, conf, nil
}
