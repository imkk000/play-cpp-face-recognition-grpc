package main

import (
	context "context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("connect to server")
	}
	defer conn.Close()

	client := NewMessageServiceClient(conn)

	ctx := context.Background()
	resp, err := client.Ping(ctx, &PingRequest{
		Message: "Hi from go client",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("ping to server")
	}

	log.Info().Msg(resp.Reply)
}
