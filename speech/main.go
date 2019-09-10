package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/storage"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"google.golang.org/grpc"
	proto "weavelab.xyz/ml-guild-speech-to-text/proto"
)

type server struct{}

func (s *server) SpeechToText(ctx context.Context, in *proto.Voice) (*proto.Text, error) {

	script, err := transcribe(ctx, in.GetVoiceLocation())
	if err != nil {
		panic(err)
	}
	uri, err := saveTranscript(ctx, in.GetBucket(), script, in.GetTextLocation(), in.GetProjectID())
	if err != nil {
		panic(err)
	}
	return &proto.Text{TestLocation: uri}, nil
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
func transcribe(ctx context.Context, filename string) (script string, err error) {
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
	var conf float32
	// Prints the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			script = fmt.Sprintf("%s%s\n", script, alt.Transcript)
			conf = conf + alt.Confidence
		}
	}
	conf = conf / float32(len(resp.Results))
	fmt.Printf("Overall - Confidence: %v\n", conf)
	return script, nil
}
func saveTranscript(ctx context.Context, bucketName, object, location, projectID string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	wc := client.Bucket(bucketName).Object(location).NewWriter(ctx)
	if _, err = io.WriteString(wc, object); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	guri := fmt.Sprintf("gs://%s/%s", bucketName, location)
	return guri, nil
}
