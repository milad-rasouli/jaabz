package jaabz

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/milad-rasouli/jaabz/internal/entity"
	"github.com/milad-rasouli/jaabz/internal/infra/godotenv"
	"log/slog"
	"strings"
)

type Jaabz struct {
	logger *slog.Logger
	url    string
	client *http.Client
}

func New(env *godotenv.Env, logger *slog.Logger) *Jaabz {
	return &Jaabz{
		logger: logger.With("repo", "jaabz"),
		url:    env.JaabzHost,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (j *Jaabz) GetJobs() ([]entity.Job, error) {
	lg := j.logger.With("method", "GetJobs")
	lg.Info("Starting to fetch jobs", "url", j.url)

	resp, err := j.client.Get(j.url)
	if err != nil {
		lg.Error("Failed to fetch URL", "error", err)
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lg.Error("Unexpected status code", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	lg.Debug("Successfully fetched HTML response")

	// Parse HTML using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		lg.Error("Failed to parse HTML", "error", err)
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	lg.Debug("Successfully parsed HTML document")

	var jobs []entity.Job

	// Find all job cards
	doc.Find("div.card-grid-2").Each(func(i int, s *goquery.Selection) {
		job := entity.Job{}

		// Extract job title
		titleSel := s.Find("h6.job-title a")
		job.Title = strings.TrimSpace(titleSel.Text())
		if href, exists := titleSel.Attr("href"); exists {
			job.VisitLink = strings.TrimSpace(href)
		}

		// Extract work status (Remote or Visa sponsorship & Relocation)
		workStatus := s.Find("strong.s-card-location").Text()
		job.WorkStatus = strings.TrimSpace(workStatus)
		if job.WorkStatus == "" {
			job.WorkStatus = "Remote" // Default to Remote if text is empty
		}

		// Extract skills
		s.Find("div.mt-20 a.btn-grey-small").Each(func(_ int, skillSel *goquery.Selection) {
			skill := strings.TrimSpace(skillSel.Text())
			if skill != "..." {
				job.Skills = append(job.Skills, skill)
			}
		})

		// Extract company name
		companySel := s.Find("div.info-right-img a")
		job.Company = strings.TrimSpace(companySel.Text())

		// Extract location
		locationSel := s.Find("span.card-location")
		job.Location = strings.TrimSpace(locationSel.Text())

		// Append job to slice
		jobs = append(jobs, job)

		lg.Debug("Extracted job details", "title", job.Title, "company", job.Company, "skills_count", len(job.Skills))
	})

	lg.Info("Completed fetching jobs", "total_jobs", len(jobs))

	return jobs, nil
}
