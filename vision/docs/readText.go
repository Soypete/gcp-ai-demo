package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
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
	filename := "testdata/IMG_1900.PNG"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}
	getBoundingBoxes(annotation)

}

func getBoundingBoxes(ann *pb.TextAnnotation) {
	pages := ann.GetPages()
	for _, p := range pages {
		blocks := p.GetBlocks()
		for i, b := range blocks {
			box := b.GetBoundingBox()
			vert := box.GetVertices()
			drawBox(vert[0].GetX(), vert[1].GetY(), vert[2].GetX(), vert[3].GetY(), i)
		}
	}
}

func drawBox(x1, y1, x2, y2 int32, count int) error {
	fSrc, err := os.Open("testdata/IMG_1900.PNG")
	if err != nil {
		log.Fatal(err)
	}
	defer fSrc.Close()
	src, err := png.Decode(fSrc)
	if err != nil {
		log.Fatal(err)
	}
	// Initialize the graphic context on an RGBA image
	dst := image.NewRGBA(image.Rect(int(x1), int(y1), int(x2), int(y2)))
	green := image.NewUniform(color.RGBA{0x00, 0x1f, 0x00, 0xff})
	draw.DrawMask(dst, src.Bounds(), src, image.Point{}, green, image.Point{}, draw.Over)

	fDst, err := os.Create(fmt.Sprintf("testdata/masks/mask-%d.png", count))
	if err != nil {
		log.Fatal(err)
	}
	defer fDst.Close()
	err = png.Encode(fDst, dst)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("dst has bounds %v.\n", dst.Bounds())

	return nil
}
