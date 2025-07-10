package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/qwaq-dev/macan-ai/pkg/pb"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
	"google.golang.org/genai"
)

func RawResumeToResumeData(rawResume string) (*pb.ResumeData, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	prompt := fmt.Sprintf(`
	Ты являешься помощником для Golang разработчика. Ты должен одать информацию о пользователе в виде такой структуры
	
	Тебе не нужно писать код на Golang. Твоя задача стоит в том, чтобы взять информацию из строки резюме и вернуть мне в виде JSON идентичного этой структуре, чтобы я могу просто анмаршалить твой ответ в готовые структуры. 
	Не больше, не меньше.

	НЕ ПИШИ НИЧЕГО, КРОМЕ САМОй JSON СТРУКТУРЫ. АБСОЛЮТНО НИЧЕГО

	type ResumeData struct {
		Success        bool                   json:"success,omitempty"
		FullName       *FullName              json:"fullName,omitempty"
		ContactInfo    *ContactInfo           json:"contactInfo,omitempty"
		Summary        string                 json:"summary,omitempty"
		Skills         []string               json:"skills,omitempty"
		WorkExperience []*WorkExperience      json:"workExperience,omitempty"
		Education      []*Education           json:"education,omitempty"
		Projects       []*PersonalProjects    json:"projects,omitempty"
		SoftSkills     []string               json:"softSkills,omitempty"
		AdditionalInfo *AdditionalInfo       
	}

	type ContactInfo struct {
		PhoneNumber     string                json:"phoneNumber,omitempty"
		Email           string                json:"email,omitempty"
		GithubUrl       string                json:"githubUrl,omitempty"
		LinkedinUrl     string                json:"linkedinUrl,omitempty"
		PersonalWebsite string                json:"personalWebsite,omitempty"
	}

	type Period struct {
		state         protoimpl.MessageState 
		Start         string                  json:"start,omitempty"
		End           string                  json:"end,omitempty"
		unknownFields protoimpl.UnknownFields
		sizeCache     protoimpl.SizeCache
	}

	type WorkExperience struct {
		CompanyName      string               json:"companyName,omitempty"
		Position         string               json:"position,omitempty"
		Period           *Period              json:"period,omitempty"
		Responsibilities []string             json:"responsibilities,omitempty"
		Technologies     []string             json:"technologies,omitempty"
	}


	type Education struct {
		Institution   string                  json:"institution,omitempty"
		Major         string                  json:"major,omitempty"
		DegreeType    string                  json:"degreeType,omitempty"
		Period        string                  json:"period,omitempty"
	}

	type PersonalProjects struct {
		ProjectName   string                  json:"projectName,omitempty"
		Period        string                  json:"period,omitempty"
		Description   string                  json:"description,omitempty"
		Technologies  []string                json:"technologies,omitempty"
	}

	type AdditionalInfo struct {
		DesiredSalary   int64                 json:"desiredSalary,omitempty"
		RelocationReady bool                  json:"relocationReady,omitempty"
		RemoteWorkReady bool                  json:"remoteWorkReady,omitempty"
	}

	Не учитывай персональный сайт, делай его null

	А вот резюме, которое ты должен превратить в эту структуру: %s 
	`, rawResume)

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error with model generation: %w", err)
	}

	var resumeData pb.ResumeData

	generatedResume := resp.Text()

	cleanedResume := strings.TrimPrefix(generatedResume, "``` json\n")
	cleanedResume = strings.TrimSuffix(cleanedResume, "\n```")
	cleanedResume = strings.TrimSpace(cleanedResume)

	log.Printf("Cleaned JSON content:\n%s\n", cleanedResume)

	if err := json.Unmarshal([]byte(cleanedResume), &resumeData); err != nil {
		fmt.Println("Error with unmarshaling", sl.Err(err))
	}

	return &resumeData, nil
}
