package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	proto "weavelab.xyz/ml-guild-speech-to-text/proto"
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

	loc, err := client.SpeechToText(ctx, &proto.Voice{VoiceLocation: "gs://phone-call-analysis/gong-calls/2019-08-02/../ml-guild-call-analysis/gong-calls/data/1131086855821845180.hi.mp3", Bucket: "phone-call-analysis", TextLocation: "voice-mails", ProjectID: os.Getenv("PROJECT_ID")})
	if err != nil {
		panic(err)
	}
	fmt.Println(loc)
}
