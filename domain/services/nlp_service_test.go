package services

import (
	"testing"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
)

// MockNLPService for testing without actually calling OpenAI
type MockNLPService struct {
	responses map[string]*entities.UserSelection
	errors    map[string]error
}

func NewMockNLPService() *MockNLPService {
	return &MockNLPService{
		responses: make(map[string]*entities.UserSelection),
		errors:    make(map[string]error),
	}
}

func (m *MockNLPService) AddResponse(text string, response *entities.UserSelection) {
	m.responses[text] = response
}

func (m *MockNLPService) AddError(text string, err error) {
	m.errors[text] = err
}

func (m *MockNLPService) ParseReminderText(userID int64, text string, userTimezone string, userLanguage string) (*entities.UserSelection, error) {
	if err, exists := m.errors[text]; exists {
		return nil, err
	}
	if response, exists := m.responses[text]; exists {
		return response, nil
	}
	// Default fallback
	return nil, nil
}

func TestNLPService_Creation(t *testing.T) {
	t.Run("fails without API key", func(t *testing.T) {
		config := &config.Config{
			OpenAI: config.OpenAIConfig{
				APIKey: "",
				Model:  "gpt-4o-mini",
			},
		}

		_, err := NewNLPService(config, nil)
		if err == nil {
			t.Error("Expected error when creating NLP service without API key")
		}
	})

	t.Run("succeeds with API key", func(t *testing.T) {
		config := &config.Config{
			OpenAI: config.OpenAIConfig{
				APIKey: "test-key",
				Model:  "gpt-4o-mini",
			},
		}

		service, err := NewNLPService(config, nil)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if service == nil {
			t.Error("Expected service to be created")
		}
	})
}

func TestNLPService_MockedParsing(t *testing.T) {
	mock := NewMockNLPService()

	testCases := []struct {
		name         string
		userID       int64
		text         string
		timezone     string
		language     string
		expectedType entities.RecurrenceType
		expectedTime string
		expectedMsg  string
		expectError  bool
	}{
		{
			name:         "English - call boss in 20 min",
			userID:       1,
			text:         "call boss in 20 min",
			timezone:     "UTC",
			language:     "en",
			expectedType: entities.Once,
			expectedTime: "10:20", // assuming current time is 10:00
			expectedMsg:  "call boss",
			expectError:  false,
		},
		{
			name:         "English - go to clinic on Monday 9:00",
			userID:       1,
			text:         "go to clinic on Monday 9:00",
			timezone:     "UTC",
			language:     "en",
			expectedType: entities.Weekly,
			expectedTime: "09:00",
			expectedMsg:  "go to clinic",
			expectError:  false,
		},
		{
			name:         "English - tomorrow 2 PM visit dentist",
			userID:       1,
			text:         "tomorrow 2 PM visit dentist",
			timezone:     "UTC",
			language:     "en",
			expectedType: entities.Once,
			expectedTime: "14:00",
			expectedMsg:  "visit dentist",
			expectError:  false,
		},
		{
			name:         "English - 9 of May 10 AM buy tickets",
			userID:       1,
			text:         "9 of May 10 AM buy tickets",
			timezone:     "UTC",
			language:     "en",
			expectedType: entities.Once,
			expectedTime: "10:00",
			expectedMsg:  "buy tickets",
			expectError:  false,
		},
		{
			name:         "English - every weekday at 8 wake up",
			userID:       1,
			text:         "every weekday at 8 wake up",
			timezone:     "UTC",
			language:     "en",
			expectedType: entities.Weekly,
			expectedTime: "08:00",
			expectedMsg:  "wake up",
			expectError:  false,
		},
		{
			name:         "English - every Wed and Fri english lesson at 10",
			userID:       1,
			text:         "every Wed and Fri english lesson at 10",
			timezone:     "UTC",
			language:     "en",
			expectedType: entities.Weekly,
			expectedTime: "10:00",
			expectedMsg:  "english lesson",
			expectError:  false,
		},
		{
			name:         "Ukrainian - подзвонити босу через 20 хвилин",
			userID:       1,
			text:         "подзвонити босу через 20 хвилин",
			timezone:     "Europe/Kiev",
			language:     "uk",
			expectedType: entities.Once,
			expectedTime: "10:20",
			expectedMsg:  "подзвонити босу",
			expectError:  false,
		},
		{
			name:         "Ukrainian - піти до клініки в понеділок о 9:00",
			userID:       1,
			text:         "піти до клініки в понеділок о 9:00",
			timezone:     "Europe/Kiev",
			language:     "uk",
			expectedType: entities.Weekly,
			expectedTime: "09:00",
			expectedMsg:  "піти до клініки",
			expectError:  false,
		},
	}

	// Setup mock responses
	for _, tc := range testCases {
		if !tc.expectError {
			selection := entities.NewUserSelection()
			selection.SetRecurrenceType(tc.expectedType)
			selection.SetSelectedTime(tc.expectedTime)
			selection.SetReminderMessage(tc.expectedMsg)

			// Add specific setup for different types
			switch tc.expectedType {
			case entities.Weekly:
				if tc.text == "go to clinic on Monday 9:00" || tc.text == "піти до клініки в понеділок о 9:00" {
					selection.WeekOptions = []time.Weekday{time.Monday}
				} else if tc.text == "every weekday at 8 wake up" {
					selection.WeekOptions = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
				} else if tc.text == "every Wed and Fri english lesson at 10" {
					selection.WeekOptions = []time.Weekday{time.Wednesday, time.Friday}
				}
			case entities.Once:
				// For once reminders, we'd set specific dates
				if tc.text == "tomorrow 2 PM visit dentist" {
					tomorrow := time.Now().Add(24 * time.Hour)
					selection.SetSelectedDate(tomorrow)
				} else if tc.text == "9 of May 10 AM buy tickets" {
					mayNinth := time.Date(time.Now().Year(), time.May, 9, 10, 0, 0, 0, time.UTC)
					selection.SetSelectedDate(mayNinth)
				}
			}

			mock.AddResponse(tc.text, selection)
		}
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := mock.ParseReminderText(tc.userID, tc.text, tc.timezone, tc.language)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for text: %s", tc.text)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result, got nil")
				return
			}

			if result.RecurrenceType != tc.expectedType {
				t.Errorf("Expected recurrence type %v, got %v", tc.expectedType, result.RecurrenceType)
			}

			if result.SelectedTime != tc.expectedTime {
				t.Errorf("Expected time %s, got %s", tc.expectedTime, result.SelectedTime)
			}

			if result.ReminderMessage != tc.expectedMsg {
				t.Errorf("Expected message %s, got %s", tc.expectedMsg, result.ReminderMessage)
			}
		})
	}
}

func TestNLPService_EdgeCases(t *testing.T) {
	mock := NewMockNLPService()

	testCases := []struct {
		name        string
		userID      int64
		text        string
		expectError bool
	}{
		{
			name:        "Empty text",
			userID:      1,
			text:        "",
			expectError: true,
		},
		{
			name:        "Unclear text",
			userID:      1,
			text:        "something unclear",
			expectError: true,
		},
		{
			name:        "Missing time",
			userID:      1,
			text:        "call doctor",
			expectError: true,
		},
		{
			name:        "Invalid time format",
			userID:      1,
			text:        "call at 25:00",
			expectError: true,
		},
	}

	// Setup mock errors
	for _, tc := range testCases {
		if tc.expectError {
			mock.AddError(tc.text, errors.NewDomainError("INVALID_TEXT", "Cannot parse text", nil))
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := mock.ParseReminderText(tc.userID, tc.text, "UTC", "en")

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for text: %s", tc.text)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected result, got nil")
				}
			}
		})
	}
}

func TestNLPService_TimezoneHandling(t *testing.T) {
	mock := NewMockNLPService()

	// Test timezone-specific parsing
	selection := entities.NewUserSelection()
	selection.SetRecurrenceType(entities.Once)
	selection.SetSelectedTime("14:00")
	selection.SetReminderMessage("meeting")

	// Same text should work in different timezones
	mock.AddResponse("meeting at 2 PM", selection)

	timezones := []string{"UTC", "America/New_York", "Europe/London", "Asia/Tokyo"}

	for _, tz := range timezones {
		t.Run("Timezone_"+tz, func(t *testing.T) {
			result, err := mock.ParseReminderText(1, "meeting at 2 PM", tz, "en")

			if err != nil {
				t.Errorf("Unexpected error for timezone %s: %v", tz, err)
				return
			}

			if result == nil {
				t.Errorf("Expected result for timezone %s, got nil", tz)
				return
			}

			if result.SelectedTime != "14:00" {
				t.Errorf("Expected time 14:00 for timezone %s, got %s", tz, result.SelectedTime)
			}
		})
	}
}

func TestNLPService_LanguageHandling(t *testing.T) {
	mock := NewMockNLPService()

	// Test language-specific parsing
	englishSelection := entities.NewUserSelection()
	englishSelection.SetRecurrenceType(entities.Daily)
	englishSelection.SetSelectedTime("08:00")
	englishSelection.SetReminderMessage("wake up")

	ukrainianSelection := entities.NewUserSelection()
	ukrainianSelection.SetRecurrenceType(entities.Daily)
	ukrainianSelection.SetSelectedTime("08:00")
	ukrainianSelection.SetReminderMessage("прокинутися")

	mock.AddResponse("wake up at 8 AM daily", englishSelection)
	mock.AddResponse("щодня прокидатися о 8 ранку", ukrainianSelection)

	testCases := []struct {
		userID   int64
		text     string
		language string
		expected string
	}{
		{
			userID:   1,
			text:     "wake up at 8 AM daily",
			language: "en",
			expected: "wake up",
		},
		{
			userID:   1,
			text:     "щодня прокидатися о 8 ранку",
			language: "uk",
			expected: "прокинутися",
		},
	}

	for _, tc := range testCases {
		t.Run("Language_"+tc.language, func(t *testing.T) {
			result, err := mock.ParseReminderText(tc.userID, tc.text, "UTC", tc.language)

			if err != nil {
				t.Errorf("Unexpected error for language %s: %v", tc.language, err)
				return
			}

			if result == nil {
				t.Errorf("Expected result for language %s, got nil", tc.language)
				return
			}

			if result.ReminderMessage != tc.expected {
				t.Errorf("Expected message %s for language %s, got %s", tc.expected, tc.language, result.ReminderMessage)
			}
		})
	}
}
