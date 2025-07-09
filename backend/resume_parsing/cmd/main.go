package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/qwaq-dev/macan-ai/pkg/pb"
	readpdf "github.com/qwaq-dev/macan-ai/resume_parsing/read_pdf"
	"google.golang.org/grpc"
)

type resumeParsingServer struct {
	pb.UnimplementedResumeParsingServiceServer
}

func (s *resumeParsingServer) SendResumePath(ctx context.Context, req *pb.ResumePath) (*pb.ResumeData, error) {
	log.Printf("Resume Parsing Service: Received request for filepath: %s", req.GetFilepath())

	resume, err := readpdf.ReadPdf(req.GetFilepath())
	if err != nil {
		log.Fatalf("Error with reading pdf:%v", err)
	}

	fmt.Println(resume)

	response := &pb.ResumeData{
		Success: true, // Указываем, что парсинг успешен (для примера)
		FullName: &pb.FullName{
			FirstName: "Иван",
			LastName:  "Иванов",
		},
		ContactInfo: &pb.ContactInfo{
			PhoneNumber:     "+77771234567",
			Email:           "ivan.ivanov@example.com",
			GithubUrl:       "github.com/ivanov",
			LinkedinUrl:     "linkedin.com/in/ivanov",
			PersonalWebsite: "ivanov.com",
		},
		Summary: "Опытный Golang разработчик с 5-летним стажем в разработке высоконагруженных систем.",
		Skills: []string{
			"Go", "gRPC", "REST API", "PostgreSQL", "Docker", "Kubernetes", "Microservices", "RabbitMQ",
		},
		WorkExperience: []*pb.WorkExperience{
			{
				CompanyName:      "Tech Solutions Inc.",
				Position:         "Senior Software Engineer",
				Period:           &pb.Period{Start: "2022-01", End: "Present"},
				Responsibilities: []string{"Разработка и поддержка микросервисов на Go", "Оптимизация производительности", "Code Review"},
				Technologies:     []string{"Go", "gRPC", "Docker", "PostgreSQL"},
			},
			{
				CompanyName:      "Innovatech LLC",
				Position:         "Software Engineer",
				Period:           &pb.Period{Start: "2019-06", End: "2021-12"},
				Responsibilities: []string{"Разработка backend API", "Интеграция сторонних сервисов"},
				Technologies:     []string{"Go", "Fiber", "MongoDB"},
			},
		},
		Education: []*pb.Education{
			{
				Institution: "Казахский Национальный Университет",
				Major:       "Информационные Системы",
				DegreeType:  "Магистр",
				Period:      "2017-2019",
			},
			{
				Institution: "Казахский Национальный Технический Университет",
				Major:       "Компьютерные Науки",
				DegreeType:  "Бакалавр",
				Period:      "2013-2017",
			},
		},
		Projects: []*pb.PersonalProjects{
			{
				ProjectName:  "E-commerce API",
				Period:       "2023-03 - 2023-09",
				Description:  "Разработка высокопроизводительного API для онлайн-магазина с использованием Go и gRPC.",
				Technologies: []string{"Go", "gRPC", "PostgreSQL", "Kafka"},
			},
		},
		SoftSkills: []string{"Коммуникабельность", "Решение проблем", "Работа в команде"},
		AdditionalInfo: &pb.AdditionalInfo{
			DesiredSalary:   350000, // Пример в тенге
			RelocationReady: true,
			RemoteWorkReady: true,
		},
	}

	log.Printf("Resume Parsing Service: Successfully processed %s, returning data.", req.GetFilepath())
	return response, nil
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
