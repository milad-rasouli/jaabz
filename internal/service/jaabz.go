package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/milad-rasouli/jaabz/internal/error_list"
	"github.com/milad-rasouli/jaabz/internal/repo/telegram"
	"strings"
	"time"

	"github.com/milad-rasouli/jaabz/internal/entity"
	"github.com/milad-rasouli/jaabz/internal/repo/duplicate"
	"github.com/milad-rasouli/jaabz/internal/repo/jaabz"
	"log/slog"
)

type JaabzService struct {
	logger        *slog.Logger
	duplicateRepo *duplicate.Duplicate
	jaabzRepo     *jaabz.Jaabz
	teleRepo      *telegram.Telegram
}

func NewJaabzService(logger *slog.Logger,
	drepo *duplicate.Duplicate,
	jrepo *jaabz.Jaabz,
	tele *telegram.Telegram) *JaabzService {
	return &JaabzService{
		logger:        logger.With("service", "jaabz"),
		duplicateRepo: drepo,
		jaabzRepo:     jrepo,
		teleRepo:      tele,
	}
}

// StartJaabzProcess starts a background process to fetch jobs every 60 seconds,
// check for duplicates, and show non-duplicate jobs.
func (j *JaabzService) StartJaabzProcess(ctx context.Context) error {
	lg := j.logger.With("method", "StartJaabzProcess")
	lg.Info("Starting Jaabz job processing")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Run immediately on start
	if err := j.processJobs(ctx); err != nil {
		lg.Error("Initial job processing failed", "error", err)
	}

	// Continue running every 60 seconds
	for {
		select {
		case <-ctx.Done():
			lg.Info("Stopping Jaabz job processing due to context cancellation")
			return ctx.Err()
		case <-ticker.C:
			if err := j.processJobs(ctx); err != nil {
				lg.Error("Job processing failed", "error", err)
			}
		}
	}
}

// processJobs fetches jobs, checks for duplicates, and shows non-duplicate jobs.
func (j *JaabzService) processJobs(ctx context.Context) error {
	lg := j.logger.With("method", "processJobs")

	lg.Debug("Fetching jobs")
	jobs, err := j.jaabzRepo.GetJobs()
	if err != nil {
		lg.Error("Failed to fetch jobs", "error", err)
		return fmt.Errorf("failed to fetch jobs: %w", err)
	}

	lg.Info("Fetched jobs", "count", len(jobs))

	for _, job := range jobs {
		err := j.duplicateRepo.SaveAndCheckDuplicate(ctx, job.VisitLink)
		if errors.Is(err, error_list.ErrDuplicate) {
			lg.Debug("Duplicate job skipped", "visit_link", job.VisitLink, "title", job.Title)
			continue
		}
		if err != nil {
			lg.Error("Failed to check duplicate", "visit_link", job.VisitLink, "error", err)
			continue
		}
		j.show(job)
		lg.Debug("Displayed non-duplicate job", "visit_link", job.VisitLink, "title", job.Title)
	}

	return nil
}

// show prints the details of the provided job.
func (j *JaabzService) show(job entity.Job) {
	fmt.Printf("Title: %s\n", job.Title)
	fmt.Printf("Company: %s\n", job.Company)
	fmt.Printf("Work Status: %s\n", job.WorkStatus)
	fmt.Printf("Visit Link: %s\n", job.VisitLink)
	fmt.Printf("Skills: %s\n", strings.Join(job.Skills, ", "))
	fmt.Printf("Location: %s\n", job.Location)
	fmt.Println("---")
}
