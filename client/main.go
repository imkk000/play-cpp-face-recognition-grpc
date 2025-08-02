package main

import (
	context "context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const inputPath = "input"

func main() {
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("connect to server")
	}
	defer conn.Close()

	client := NewFaceRecognitionServiceClient(conn)

	files, _ := os.ReadDir(inputPath)
	images := make([][]byte, len(files))
	for i, file := range files {
		fmt.Println("Reading file:", file.Name())
		images[i], _ = os.ReadFile(filepath.Join(inputPath, file.Name()))
	}

	ctx := context.Background()
	resp, err := client.DetectFaces(ctx, &DetectFacesRequest{Images: images})
	if err != nil {
		log.Fatal().Err(err).Msg("ping to server")
	}

	log.Info().Msgf("status: %v", resp.GetStatus())

	for i, face := range resp.GetFaces() {
		log.Info().Msgf("%d x: %f, y: %f", i, face.GetX(), face.GetY())
		log.Info().Msgf("%d w: %f, h: %f", i, face.GetWidth(), face.GetHeight())
		log.Info().Msgf("%d confidence: %f", i, face.GetConfidence())
		log.Info().Msgf("%d landmarks: %v", i, face.GetLandmarks())
	}
}
