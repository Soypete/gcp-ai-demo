package main

import (
	"context"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// sets the name of the image file to annotate.
	filename := "testdata/max-the-borkie.jpg"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}
	client.DetectDocumentText(ctx, image, nil)

}
