package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/sashabaranov/go-openai"
)

// NLPError represents different types of NLP-related errors
type NLPError struct {
	Type    string
	Message string
	Code    string
}

func (e *NLPError) Error() string {
	return e.Message
}

// Error types
const (
	NLPErrorRateLimit = "RATE_LIMIT_EXCEEDED"
	NLPErrorParsing   = "PARSING_FAILED"
	NLPErrorInternal  = "INTERNAL_ERROR"
)

// NLPService handles natural language processing for reminder creation
type NLPService interface {
	ParseReminderText(userID int64, text string, userTimezone string, userLanguage string) (*entities.UserSelection, error)
}

type nlpService struct {
	client    *openai.Client
	config    *config.Config
	usageRepo repositories.PremiumUsageRepository
}

// ReminderRequest represents the structure we want OpenAI to return
type ReminderRequest struct {
	RecurrenceType  string         `json:"recurrenceType"`
	WeekOptions     []time.Weekday `json:"weekOptions,omitempty"`
	MonthOptions    []int          `json:"monthOptions,omitempty"`
	SelectedDate    string         `json:"selectedDate,omitempty"` // ISO format date
	SelectedTime    string         `json:"selectedTime"`           // HH:MM format
	IntervalDays    int            `json:"intervalDays,omitempty"`
	ReminderMessage string         `json:"reminderMessage"`
	IsValid         bool           `json:"isValid"`
	ErrorMessage    string         `json:"errorMessage,omitempty"`
}

// NewNLPService creates a new NLP service
func NewNLPService(config *config.Config, usageRepo repositories.PremiumUsageRepository) (NLPService, error) {
	if config.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(config.OpenAI.APIKey)
	return &nlpService{
		client:    client,
		config:    config,
		usageRepo: usageRepo,
	}, nil
}

// ParseReminderText uses OpenAI to parse natural language text into a UserSelection
func (s *nlpService) ParseReminderText(userID int64, text string, userTimezone string, userLanguage string) (*entities.UserSelection, error) {
	premiumUsageResult := s.validatePremiumUsage(userID)
	if premiumUsageResult != nil {
		return nil, premiumUsageResult
	}

	prompt := s.buildPrompt(text, userTimezone, userLanguage)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.config.OpenAI.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: s.getSystemPrompt(userTimezone),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1,
			MaxTokens:   500,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	content = strings.TrimSpace(content)

	// Remove markdown code blocks if present
	if after, ok := strings.CutPrefix(content, "```json"); ok {
		content = after
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	} else if after0, ok0 := strings.CutPrefix(content, "```"); ok0 {
		content = after0
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	var reminderReq ReminderRequest
	if err := json.Unmarshal([]byte(content), &reminderReq); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w\nResponse: %s", err, content)
	}

	if !reminderReq.IsValid {
		return nil, fmt.Errorf("incomplete or invalid request: %s", reminderReq.ErrorMessage)
	}

	return s.convertToUserSelection(&reminderReq, userTimezone)
}

func (s *nlpService) validatePremiumUsage(userID int64) error {
	// Get or create user usage record
	usage, err := s.usageRepo.GetOrCreateUserUsage(userID)
	if err != nil {
		log.Printf("Failed to get NLP usage for user %d: %v", userID, err)
		return &NLPError{
			Type:    NLPErrorInternal,
			Message: "Failed to check usage limits",
			Code:    "USAGE_CHECK_FAILED",
		}
	}

	// Check if user can make a request
	if usage.CanMakeRequest() {
		return nil
	}

	log.Printf("User %d exceeded NLP rate limit: %d/%d requests used", userID, usage.RequestsUsed, usage.RequestsLimit)

	remainingDays := s.getDaysUntilReset(usage)
	var errorMsg string

	switch usage.PremiumStatus {
	case entities.PremiumStatusFree:
		errorMsg = fmt.Sprintf("You've reached your monthly limit of %d AI text reminders. Upgrade to Premium for more requests or try again in %d days.",
			usage.RequestsLimit, remainingDays)
	case entities.PremiumStatusBasic:
		errorMsg = fmt.Sprintf("You've reached your monthly limit of %d AI text reminders. Upgrade to Pro for unlimited requests or try again in %d days.",
			usage.RequestsLimit, remainingDays)
	default:
		errorMsg = "Rate limit exceeded. Please try again later."
	}

	return &NLPError{
		Type:    NLPErrorRateLimit,
		Message: errorMsg,
		Code:    "MONTHLY_LIMIT_EXCEEDED",
	}
}

// getDaysUntilReset calculates days until the monthly reset
func (s *nlpService) getDaysUntilReset(usage *entities.PremiumUsage) int {
	now := usage.LastReset
	nextMonth := now.AddDate(0, 1, 0)

	// Get the first day of next month
	firstOfNextMonth := nextMonth.AddDate(0, 0, -nextMonth.Day()+1)

	days := int(firstOfNextMonth.Sub(now).Hours() / 24)
	if days < 1 {
		days = 1
	}
	return days
}

// getSystemPrompt returns the system prompt for OpenAI
func (s *nlpService) getSystemPrompt(userTimezone string) string {
	return fmt.Sprintf(`You are a helpful assistant that converts natural language reminder requests into structured JSON format. 

Current timezone: %s
Current date and time: %s

You must respond ONLY with valid JSON in the following format:
{
    "recurrenceType": "Once|Daily|Weekly|Monthly|Interval",
    "weekOptions": [0,1,2,3,4,5,6], // Only for Weekly - Sunday=0, Monday=1, etc.
    "monthOptions": [1,2,3,...,31], // Only for Monthly - days of month
    "selectedDate": "2025-01-15", // ISO date format, only for Once
    "selectedTime": "14:30", // HH:MM format (24-hour)
    "intervalDays": 5, // Only for Interval type
    "reminderMessage": "extracted message",
    "isValid": true, // false if request is incomplete or unclear
    "errorMessage": "reason why invalid" // only if isValid is false
}

Rules:
1. For "Once": set selectedDate and selectedTime
2. For "Daily": set selectedTime only
3. For "Weekly": set selectedTime and weekOptions (array of weekday numbers)
4. For "Monthly": set selectedTime and monthOptions (array of day numbers)
5. For "Interval": set selectedTime and intervalDays
6. Time parsing: "in X minutes/hours" means from now, "at X" means specific time, "tomorrow" means next day
7. Week parsing: "weekdays" = [1,2,3,4,5], "weekends" = [0,6], "every day" = Daily
8. If time is missing or unclear, set isValid to false
9. Extract the actual reminder message/task from the text
10. Handle both English and Ukrainian text`, userTimezone, time.Now().Format("2006-01-02 15:04:05 MST"))
}

// buildPrompt creates the user prompt
func (s *nlpService) buildPrompt(text string, userTimezone string, userLanguage string) string {
	return fmt.Sprintf("Parse this reminder request: \"%s\"\nUser timezone: %s\nUser language: %s",
		text, userTimezone, userLanguage)
}

// convertToUserSelection converts ReminderRequest to UserSelection
func (s *nlpService) convertToUserSelection(req *ReminderRequest, userTimezone string) (*entities.UserSelection, error) {
	selection := entities.NewUserSelection()

	// Set recurrence type
	recurrenceType, err := entities.ToRecurrenceType(req.RecurrenceType)
	if err != nil {
		return nil, fmt.Errorf("invalid recurrence type: %s", req.RecurrenceType)
	}
	selection.SetRecurrenceType(recurrenceType)

	// Set time
	selection.SetSelectedTime(req.SelectedTime)

	// Set message
	selection.SetReminderMessage(req.ReminderMessage)

	// Set type-specific options
	switch recurrenceType {
	case entities.Once:
		if req.SelectedDate != "" {
			date, err := time.Parse("2006-01-02", req.SelectedDate)
			if err != nil {
				return nil, fmt.Errorf("invalid date format: %s", req.SelectedDate)
			}

			// Parse the time and combine with date
			timeParts := strings.Split(req.SelectedTime, ":")
			if len(timeParts) != 2 {
				return nil, fmt.Errorf("invalid time format: %s", req.SelectedTime)
			}

			// Load user timezone
			loc, err := time.LoadLocation(userTimezone)
			if err != nil {
				loc = time.UTC
			}

			// Create datetime in user's timezone
			var hour, minute int
			fmt.Sscanf(req.SelectedTime, "%d:%d", &hour, &minute)
			dateTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc)
			selection.SetSelectedDate(dateTime)
		}
	case entities.Weekly:
		selection.WeekOptions = req.WeekOptions
	case entities.Monthly:
		selection.MonthOptions = req.MonthOptions
	case entities.Interval:
		selection.SetCustomInterval(req.IntervalDays)
	}

	return selection, nil
}
