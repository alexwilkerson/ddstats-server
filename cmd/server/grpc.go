package main

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/alexwilkerson/ddstats-server/gamesubmission"
	"github.com/alexwilkerson/ddstats-server/pkg/models/postgres"
)

const (
	oldestValidClientVersion = "0.6.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGameRecorderServer
	db                   *postgres.Postgres
	client               *http.Client
	currentClientVersion string
}

// SubmitGame is uncommented so far.
func (s *server) SubmitGame(ctx context.Context, in *pb.SubmitGameRequest) (*pb.SubmitGameReply, error) {
	duplicate, gameID, err := s.db.GameSubmissions.CheckDuplicate(in)
	if err != nil {
		return nil, fmt.Errorf("SubmitGame: error checking for duplicate game: %w", err)
	}

	if duplicate {
		return &pb.SubmitGameReply{GameID: gameID}, nil
	}

	gameID, err = s.db.GameSubmissions.Insert(in)
	if err != nil {
		return nil, fmt.Errorf("SubmitGame: error inserting game: %w", err)
	}

	return &pb.SubmitGameReply{GameID: gameID}, nil
}

func (s *server) ClientStart(ctx context.Context, in *pb.ClientStartRequest) (*pb.ClientStartReply, error) {
	motd, err := s.db.MOTD.Get()
	if err != nil {
		return nil, err
	}

	valid, err := validVersion(in.GetVersion())
	if err != nil {
		return nil, err
	}
	update, err := s.updateAvailable(in.GetVersion())
	if err != nil {
		return nil, err
	}

	return &pb.ClientStartReply{
		Motd:            motd.Message,
		ValidVersion:    valid,
		UpdateAvailable: update,
	}, nil
}

func validVersion(version string) (bool, error) {
	var vMajor, vMinor, vPatch, ovMajor, ovMinor, ovPatch int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &vMajor, &vMinor, &vPatch)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(oldestValidClientVersion, "%d.%d.%d", &ovMajor, &ovMinor, &ovPatch)
	if err != nil {
		return false, err
	}
	if vMajor > ovMajor ||
		(vMajor == ovMajor && vMinor > ovMinor) ||
		(vMajor == ovMajor && vMinor == ovMinor && vPatch >= ovPatch) {
		return true, nil
	}
	return false, nil
}

func (s *server) updateAvailable(version string) (bool, error) {
	var vMajor, vMinor, vPatch, cvMajor, cvMinor, cvPatch int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &vMajor, &vMinor, &vPatch)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(s.currentClientVersion, "%d.%d.%d", &cvMajor, &cvMinor, &cvPatch)
	if err != nil {
		return false, err
	}
	if cvMajor > vMajor ||
		(cvMajor == vMajor && cvMinor > vMinor) ||
		(cvMajor == vMajor && cvMinor == vMinor && cvPatch > vPatch) {
		return true, nil
	}
	return false, nil
}
