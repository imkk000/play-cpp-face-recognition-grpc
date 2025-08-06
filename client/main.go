package main

import (
	context "context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inputPath   = "input"
	profilePath = "profile"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("connect to server")
	}
	defer conn.Close()

	client := NewFaceRecognitionServiceClient(conn)

	rootCmd := cli.Command{
		Commands: []*cli.Command{
			{
				Name: "batch",
				Action: func(_ context.Context, c *cli.Command) error {
					files, _ := os.ReadDir(inputPath)
					images := make([][]byte, 0, len(files))
					for _, file := range files {
						ext := filepath.Ext(strings.ToLower(file.Name()))
						if file.IsDir() || ext != ".jpg" && ext != ".jpeg" {
							continue
						}
						fmt.Println("Reading file:", file.Name())
						image, _ := os.ReadFile(filepath.Join(inputPath, file.Name()))
						images = append(images, image)
					}
					images = slices.Clip(images)

					ctx := context.Background()
					resp, err := client.DetectFaces(ctx, &DetectFacesRequest{Images: images})
					if err != nil {
						return fmt.Errorf("detect faces: %w", err)
					}
					landmarksData := make([]*RecognizeFacesRequest_Face, len(resp.GetFaces()))
					log.Info().Msgf("status: %v", resp.GetStatus())
					for i, face := range resp.GetFaces() {
						log.Info().Msgf("%d x: %f, y: %f", i, face.GetX(), face.GetY())
						log.Info().Msgf("%d w: %f, h: %f", i, face.GetWidth(), face.GetHeight())
						log.Info().Msgf("%d confidence: %f", i, face.GetConfidence())
						log.Info().Msgf("%d landmarks: %v", i, face.GetLandmarks())

						body := []float32{face.GetX(), face.GetY(), face.GetWidth(), face.GetHeight()}
						body = append(body, face.GetLandmarks()...)
						body = append(body, face.GetConfidence())

						landmarksData[i] = &RecognizeFacesRequest_Face{Landmarks: body}
					}

					files, _ = os.ReadDir(profilePath)
					images = make([][]byte, 0, len(files))
					filenames := make([]string, 0, len(files))
					for _, file := range files {
						ext := filepath.Ext(strings.ToLower(file.Name()))
						if file.IsDir() || ext != ".jpg" && ext != ".jpeg" {
							continue
						}
						fmt.Println("Reading file:", file.Name())
						image, _ := os.ReadFile(filepath.Join(profilePath, file.Name()))
						images = append(images, image)
						filenames = append(filenames, file.Name())
					}
					images = slices.Clip(images)
					filenames = slices.Clip(filenames)

					for i, image := range images {
						resp, err := client.RecognizeFaces(ctx, &RecognizeFacesRequest{
							Faces: landmarksData,
							Image: image,
						})
						if err != nil {
							return fmt.Errorf("recognize faces: %s %w", filenames[i], err)
						}
						log.Info().Msgf("index: %d: filename: %s, status: %v, matched: %v", i, filenames[i], resp.GetStatus(), resp.GetValid())
					}

					return nil
				},
			},
			{
				Name: "validate",
				Action: func(_ context.Context, c *cli.Command) error {
					files, _ := os.ReadDir(inputPath)
					images := make([][]byte, len(files))
					for i, file := range files {
						fmt.Println("Reading file:", file.Name())
						images[i], _ = os.ReadFile(filepath.Join(inputPath, file.Name()))
					}

					ctx := context.Background()
					resp, err := client.DetectFaces(ctx, &DetectFacesRequest{Images: images})
					if err != nil {
						return fmt.Errorf("detect faces: %w", err)
					}
					log.Info().Msgf("status: %v", resp.GetStatus())
					for i, face := range resp.GetFaces() {
						log.Info().Msgf("%d x: %f, y: %f", i, face.GetX(), face.GetY())
						log.Info().Msgf("%d w: %f, h: %f", i, face.GetWidth(), face.GetHeight())
						log.Info().Msgf("%d confidence: %f", i, face.GetConfidence())
						log.Info().Msgf("%d landmarks: %v", i, face.GetLandmarks())
					}
					return nil
				},
			},
			{
				Name: "verify",
				Action: func(_ context.Context, c *cli.Command) error {
					files, _ := os.ReadDir(profilePath)
					images := make([][]byte, len(files))
					for i, file := range files {
						fmt.Println("Reading file:", file.Name())
						images[i], _ = os.ReadFile(filepath.Join(profilePath, file.Name()))
					}

					ctx := context.Background()
					resp, err := client.RecognizeFaces(ctx, &RecognizeFacesRequest{
						Faces: []*RecognizeFacesRequest_Face{
							{
								Landmarks: []float32{},
							},
							{
								Landmarks: []float32{},
							},
							{
								Landmarks: []float32{},
							},
						},
						Image: images[0],
					})
					if err != nil {
						return fmt.Errorf("recognize faces: %w", err)
					}
					log.Info().Msgf("status: %v", resp.GetStatus())
					log.Info().Msgf("matched: %v", resp.GetValid())
					return nil
				},
			},
		},
	}
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("run command")
	}
}
