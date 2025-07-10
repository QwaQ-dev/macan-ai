package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/qwaq-dev/macan-ai/pkg/pb"
	"github.com/qwaq-dev/macan-ai/resume_parsing/ai"
	readpdf "github.com/qwaq-dev/macan-ai/resume_parsing/read_pdf"
	"google.golang.org/grpc"
)

type resumeParsingServer struct {
	pb.UnimplementedResumeParsingServiceServer
}

func (s *resumeParsingServer) SendResumePath(ctx context.Context, req *pb.ResumePath) (*pb.ResumeData, error) {
	log.Printf("Resume Parsing Service: Received request for filepath: %s", req.GetFilepath())

	rawResume, err := readpdf.ReadPdf(req.GetFilepath())
	if err != nil {
		log.Fatalf("Error with reading pdf:%v", err)
	}

	resume, err := ai.RawResumeToResumeData(rawResume)
	if err != nil {
		return nil, fmt.Errorf("Error with model:%v", err)
	}

	log.Printf("Resume Parsing Service: Successfully processed %s, returning data.", req.GetFilepath())
	return resume, nil
}

func main() {
	port := "50051"
	addr := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Resume Parsing Service: Failed to listen on %s: %v", addr, err)
	}

	s := grpc.NewServer()
	pb.RegisterResumeParsingServiceServer(s, &resumeParsingServer{})

	log.Printf("Resume Parsing Service: gRPC server listening on %s", addr)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Resume Parsing Service: Failed to serve gRPC: %v", err)
	}
}
