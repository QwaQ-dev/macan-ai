package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/qwaq-dev/macan-ai/pkg/pb"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ResumeParsingClient struct {
	Client pb.ResumeParsingServiceClient
	Conn   *grpc.ClientConn
	log    *slog.Logger
}

func NewResumeParsingClient(log *slog.Logger, address string) (*ResumeParsingClient, error) {
	const op = "services.NewResumeParsingClient"
	log = log.With(slog.String("op", op))

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(address, dialOptions...)
	if err != nil {
		log.Error("Failed to connect to Resume Parsing Service", sl.Err(err), slog.String("address", address))
		return nil, fmt.Errorf("failed to connect to Resume Parsing Service at %s: %w", address, err)
	}

	client := pb.NewResumeParsingServiceClient(conn)

	log.Info("Successfully created gRPC client for Resume Parsing Service", slog.String("address", address))

	return &ResumeParsingClient{
		Client: client,
		Conn:   conn,
		log:    log,
	}, nil
}

func (r *ResumeParsingClient) Close() error {
	if r.Conn != nil {
		r.log.Debug("Closing gRPC client connection to Resume Parsing Service...")
		return r.Conn.Close()
	}
	return nil
}

func (r *ResumeParsingClient) ParseResume(filePath string) (*pb.ResumeData, error) {
	const op = "services.ResumeParsingClient.ParseResume"
	log := r.log.With("op", op)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Устанавливаем таймаут для RPC-вызова
	defer cancel()

	req := &pb.ResumePath{Filepath: filePath}
	log.Debug("Sending resume path to Resume Parsing Service", slog.String("filepath", filePath))

	resp, err := r.Client.SendResumePath(ctx, req)
	if err != nil {
		log.Error("Failed to call SendResumePath", sl.Err(err))
		return nil, fmt.Errorf("failed to call SendResumePath: %w", err)
	}

	log.Debug("Received response from Resume Parsing Service", slog.Bool("success", resp.GetSuccess()))
	// Можно добавить проверку на !resp.GetSuccess() здесь и вернуть соответствующую ошибку
	if !resp.GetSuccess() {
		log.Warn("Resume parsing service reported failure (internal service error)", slog.String("filepath", filePath))
		return resp, fmt.Errorf("resume parsing service reported non-successful status for file: %s", filePath)
	}

	return resp, nil
}
