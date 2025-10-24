package services

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// MockOpenAIClient is a mock implementation of OpenAI client for testing
type MockOpenAIClient struct {
	responses []openai.ChatCompletionResponse
	current   int
}

// NewMockOpenAIClient creates a new mock OpenAI client with predefined responses
func NewMockOpenAIClient() *MockOpenAIClient {
	// Hardcoded responses for different test scenarios
	responses := []openai.ChatCompletionResponse{
		// Response 1: Simple once reminder
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: `{
							"recurrenceType": "Once",
							"selectedDate": "2025-01-15",
							"selectedTime": "14:30",
							"reminderMessage": "Buy groceries",
							"isValid": true
						}`,
					},
				},
			},
		},
		// Response 2: Daily reminder
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: `{
							"recurrenceType": "Daily",
							"selectedTime": "09:00",
							"reminderMessage": "Take vitamins",
							"isValid": true
						}`,
					},
				},
			},
		},
		// Response 3: Weekly reminder
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: `{
							"recurrenceType": "Weekly",
							"weekOptions": [1, 3, 5],
							"selectedTime": "10:00",
							"reminderMessage": "Go to gym",
							"isValid": true
						}`,
					},
				},
			},
		},
		// Response 4: Monthly reminder
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: `{
							"recurrenceType": "Monthly",
							"monthOptions": [1, 15],
							"selectedTime": "12:00",
							"reminderMessage": "Pay bills",
							"isValid": true
						}`,
					},
				},
			},
		},
		// Response 5: Interval reminder
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: `{
							"recurrenceType": "Interval",
							"intervalDays": 3,
							"selectedTime": "16:00",
							"reminderMessage": "Water plants",
							"isValid": true
						}`,
					},
				},
			},
		},
		// Response 6: Invalid reminder (missing time)
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: `{
							"recurrenceType": "Once",
							"reminderMessage": "Call mom",
							"isValid": false,
							"errorMessage": "Time not specified"
						}`,
					},
				},
			},
		},
	}

	return &MockOpenAIClient{
		responses: responses,
		current:   0,
	}
}

// CreateChatCompletion mocks the OpenAI chat completion API
func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	// Return responses in rotation
	response := m.responses[m.current%len(m.responses)]
	m.current++

	return response, nil
}

// SetResponse allows setting a custom response for testing specific scenarios
func (m *MockOpenAIClient) SetResponse(jsonResponse string) {
	m.responses = []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: jsonResponse,
					},
				},
			},
		},
	}
	m.current = 0
}

// SetResponses allows setting multiple custom responses
func (m *MockOpenAIClient) SetResponses(jsonResponses []string) {
	responses := make([]openai.ChatCompletionResponse, len(jsonResponses))
	for i, jsonResp := range jsonResponses {
		responses[i] = openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: jsonResp,
					},
				},
			},
		}
	}
	m.responses = responses
	m.current = 0
}

// GetCallCount returns the number of times CreateChatCompletion was called
func (m *MockOpenAIClient) GetCallCount() int {
	return m.current
}

// Reset resets the mock to its initial state
func (m *MockOpenAIClient) Reset() {
	m.current = 0
}

// OpenAIClientInterface defines the interface that both real and mock clients implement
type OpenAIClientInterface interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// Verify that both implementations satisfy the interface
var _ OpenAIClientInterface = (*openai.Client)(nil)
var _ OpenAIClientInterface = (*MockOpenAIClient)(nil)
