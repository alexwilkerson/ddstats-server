package main

import (
	"context"
	"net/http"

	pb "github.com/alexwilkerson/ddstats-server/gamesubmission"
	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGameRecorderServer
	db     *postgres.Postgres
	client *http.Client
}

// SubmitGame is uncommented so far.
func (s *server) SubmitGame(ctx context.Context, in *pb.SubmitGameRequest) (*pb.SubmitGameReply, error) {
	gameID, err := s.db.GameSubmissions.Insert(in)
	if err != nil {
		return nil, err
	}
	// log.Printf("Received: %v", in.GetName())
	return &pb.SubmitGameReply{GameID: gameID}, nil
}
