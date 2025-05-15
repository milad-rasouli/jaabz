package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/milad-rasouli/jaabz/internal/entity"
	"github.com/milad-rasouli/jaabz/internal/infra/godotenv"
	"log/slog"
	"strings"
	"time"
)

type Telegram struct {
	logger    *slog.Logger
	bot       *tgbotapi.BotAPI
	channelID string
}

func New(logger *slog.Logger, env *godotenv.Env) (*Telegram, error) {
	logger = logger.With("package", "telegram")

	// Initialize Telegram bot
	bot, err := tgbotapi.NewBotAPI(env.TelegramBotToken)
	if err != nil {
		logger.Error("Failed to initialize Telegram bot", "error", err)
		return nil, fmt.Errorf("failed to initialize telegram bot: %w", err)
	}
	logger.Info("Telegram bot initialized", "bot_username", bot.Self.UserName)

	return &Telegram{
		logger:    logger.With("repo", "telegram"),
		bot:       bot,
		channelID: env.TelegramChannelID,
	}, nil
}

func (t *Telegram) Ready() bool {
	return t.bot != nil
}
func (t *Telegram) Post(job entity.Job) error {
	lg := t.logger.With("method", "Post", "job_title", job.Title, "visit_link", job.VisitLink)

	if t.bot == nil {
		lg.Error("Telegram bot not initialized")
		return fmt.Errorf("telegram bot not initialized")
	}

	message := fmt.Sprintf(
		"*New Job Posting*\n\n"+
			"*Title*: %s\n"+
			"*Company*: %s\n"+
			"*Work Status*: %s\n"+
			"*Location*: %s\n"+
			"*Skills*: %s\n"+
			"*Apply*: [Link](%s)",
		escapeMarkdown(job.Title),
		escapeMarkdown(job.Company),
		escapeMarkdown(job.WorkStatus),
		escapeMarkdown(job.Location),
		escapeMarkdown(strings.Join(job.Skills, ", ")),
		job.VisitLink,
	)

	msg := tgbotapi.NewMessageToChannel(t.channelID, message)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.DisableWebPagePreview = true

	_, err := t.bot.Send(msg)
	if err != nil {
		// Check for rate limit response
		if strings.Contains(err.Error(), "Too Many Requests") {
			var waitSec int
			_, scanErr := fmt.Sscanf(err.Error(), "Too Many Requests: retry after %d", &waitSec)
			if scanErr == nil {
				lg.Warn("Rate limited, retrying after delay", "wait_seconds", waitSec)
				time.Sleep(time.Duration(waitSec+1) * time.Second) // Add +1 for safety
				_, retryErr := t.bot.Send(msg)
				if retryErr != nil {
					lg.Error("Retry failed", "error", retryErr)
					return fmt.Errorf("telegram retry failed: %w", retryErr)
				}
				lg.Info("Successfully posted job after retry")
				return nil
			}
		}

		lg.Error("Failed to post job to Telegram channel", "error", err)
		return fmt.Errorf("failed to post job to Telegram: %w", err)
	}

	lg.Info("Successfully posted job to Telegram channel")
	return nil
}

// escapeMarkdown escapes special characters for Telegram MarkdownV2
func escapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}
