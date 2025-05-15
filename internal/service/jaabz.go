package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/milad-rasouli/jaabz/internal/error_list"
	"github.com/milad-rasouli/jaabz/internal/repo/duplicate"
	"github.com/milad-rasouli/jaabz/internal/repo/jaabz"
	"github.com/milad-rasouli/jaabz/internal/repo/telegram"
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
// check for duplicates, and post non-duplicate jobs to Telegram.
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

// processJobs fetches jobs, checks for duplicates, and posts non-duplicate jobs to Telegram.
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
		if err := j.teleRepo.Post(job); err != nil {
			lg.Error("Failed to post job to Telegram", "visit_link", job.VisitLink, "title", job.Title, "error", err)
			continue
		}
		lg.Debug("Posted non-duplicate job to Telegram", "visit_link", job.VisitLink, "title", job.Title)
	}

	return nil
}
