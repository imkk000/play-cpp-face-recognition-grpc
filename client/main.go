package main

import (
	context "context"
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	inputPath   = "input"
	profilePath = "profile"
)

func main() {
	conn, err := grpc.NewClient("stg-face-recognition-api-internal.rizzup.com:8443", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	// conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
					landmarksData := make([]*Face, len(resp.GetFaces()))
					log.Info().Msgf("status: %v", resp.GetStatus())
					for i, face := range resp.GetFaces() {
						log.Info().Msgf("%d landmarks: %v", i, face.GetLandmarks())

						landmarksData[i] = &Face{Landmarks: face.GetLandmarks()}
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

					for {
						for i, image := range images {
							start := time.Now()
							resp, err := client.RecognizeFaces(ctx, &RecognizeFacesRequest{
								Faces: landmarksData,
								Image: image,
							})
							if err != nil {
								err = fmt.Errorf("recognize faces: %s %w", filenames[i], err)
								log.Error().Msgf("index: %d: filename: %s, status: %v, since: %s", i, filenames[i], err, time.Since(start))
								continue
							}
							log.Info().Msgf("index: %d: filename: %s, status: %v, matched: %v, since: %s", i, filenames[i], resp.GetStatus(), resp.GetValid(), time.Since(start))
						}
					}

					return nil
				},
			},
		},
	}
	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("run command")
	}
}
