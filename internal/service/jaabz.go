package service

import (
	"fmt"
	"github.com/milad-rasouli/jaabz/internal/repo/duplicate"
	"github.com/milad-rasouli/jaabz/internal/repo/jaabz"
	"log/slog"
	"strings"
)

type JaabzService struct {
	logger        *slog.Logger
	duplicateRepo *duplicate.Duplicate
	jaabzRepo     *jaabz.Jaabz
}

func NewJaabzService(logger *slog.Logger,
	drepo *duplicate.Duplicate,
	jrepo *jaabz.Jaabz) *JaabzService {
	return &JaabzService{
		logger:        logger.With("service", "jaabz"),
		duplicateRepo: drepo,
		jaabzRepo:     jrepo,
	}
}

func (j *JaabzService) JaabzProcess() error {
	lg := j.logger.With("method", "JaabzProcess")
	jobs, err := j.jaabzRepo.GetJobs()
	if err != nil {
		lg.Error("Error scraping jobs: %v", err)
	}

	// Print the extracted jobs
	for i, job := range jobs {
		fmt.Printf("Job %d:\n", i+1)
		fmt.Printf("Title: %s\n", job.Title)
		fmt.Printf("Company: %s\n", job.Company)
		fmt.Printf("Work Status: %s\n", job.WorkStatus)
		fmt.Printf("Visit Link: %s\n", job.VisitLink)
		fmt.Printf("Skills: %s\n", strings.Join(job.Skills, ", "))
		fmt.Printf("Location: %s\n", job.Location)
		fmt.Println("---")
	}

	return nil
}
