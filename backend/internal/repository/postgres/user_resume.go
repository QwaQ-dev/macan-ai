package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lib/pq"
	"github.com/qwaq-dev/macan-ai/pkg/pb"
	"github.com/qwaq-dev/macan-ai/pkg/sl"
)

type UserResumeRepo struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserResumeRepo(log *slog.Logger, db *sql.DB) *UserResumeRepo {
	return &UserResumeRepo{
		log: log,
		db:  db,
	}
}

func (u *UserResumeRepo) AddResume(resume *pb.ResumeData, userId int) error {
	const op = "postgres.user_resume.AddResume"
	log := u.log.With("op", op)

	query := `
	INSERT INTO user_resume (full_name, contact_info, summary, key_skills, 
								work_experience, education, personal_projects, soft_skills, additional_info, user_id)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	log.Info("resume", resume)
	fullName := fmt.Sprintf("%s %s", resume.FullName.FirstName+resume.FullName.LastName)

	contactInfoJSON, err := json.Marshal(resume.GetContactInfo())
	if err != nil {
		log.Error("Error with marshaling contact info", sl.Err(err))
		return err
	}

	summary := resume.GetSummary()

	keySkills := pq.Array(resume.GetSkills())

	var workExperienceArray []string
	for _, exp := range resume.GetWorkExperience() {
		expJSON, err := json.Marshal(exp)
		if err != nil {
			log.Error("Error with marshaling work exp", sl.Err(err))
			return err
		}

		workExperienceArray = append(workExperienceArray, string(expJSON))
	}
	workExperience := pq.Array(workExperienceArray)

	var educationArray []string
	for _, edc := range resume.GetEducation() {
		eduJSON, err := json.Marshal(edc)
		if err != nil {
			log.Error("Error with marhaling education info", sl.Err(err))
			return err
		}
		educationArray = append(educationArray, string(eduJSON))
	}
	education := pq.Array(educationArray)

	var personalProjectsArray []string
	for _, proj := range resume.GetProjects() {
		projJSON, err := json.Marshal(proj)
		if err != nil {
			log.Error("Error with marshaling personal projects info", sl.Err(err))
			return err
		}
		personalProjectsArray = append(personalProjectsArray, string(projJSON))
	}
	personalProjects := pq.Array(personalProjectsArray)

	softSkills := pq.Array(resume.GetSoftSkills())

	var additionalInfoJSON []byte
	if resume.GetAdditionalInfo() != nil {
		additionalInfoJSON, err = json.Marshal(resume.GetAdditionalInfo())
		if err != nil {
			log.Error("Error with marshaling additional info", sl.Err(err))
			return err
		}
	} else {
		additionalInfoJSON = []byte("null")
	}

	_, err = u.db.Exec(
		query,
		fullName,
		contactInfoJSON,
		summary,
		keySkills,
		workExperience,
		education,
		personalProjects,
		softSkills,
		additionalInfoJSON, userId,
	)

	if err != nil {
		log.Error("Error with inserting data to db", sl.Err(err))
		return err
	}

	return nil
}
