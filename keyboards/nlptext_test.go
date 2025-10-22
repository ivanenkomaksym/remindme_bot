package keyboards

import (
	"errors"
	"testing"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// Mock NLP service for testing
type mockNLPService struct {
	shouldFail bool
	result     *entities.UserSelection
}

func (m *mockNLPService) ParseReminderText(userID int64, text, timezone, language string) (*entities.UserSelection, error) {
	if m.shouldFail {
		return nil, errors.New("mock NLP parsing failed")
	}
	return m.result, nil
}

func TestHandleNlpTextInputCallback(t *testing.T) {
	// Create a test user
	user := &entities.User{
		ID:       123,
		Language: LangEN,
	}

	// Mock function for updating user selection
	var capturedSelection *entities.UserSelection
	mockUpdateUserSelection := func(userID int64, selection *entities.UserSelection) error {
		if userID != user.ID {
			t.Errorf("Expected userID %d, got %d", user.ID, userID)
		}
		capturedSelection = selection
		return nil
	}

	// Call the function
	result, err := HandleNlpTextInputCallback(user, mockUpdateUserSelection)

	// Verify no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Verify the text contains expected elements
	s := T(user.Language)
	expectedText := s.NlpMenuTitle + "\n\n" + s.NlpInstructions + "\n\n" + s.NlpExamples + "\n\n" + s.NlpEnterText
	if result.Text != expectedText {
		t.Errorf("Expected text to contain NLP instructions, got: %s", result.Text)
	}

	// Verify markup has back button
	if result.Markup == nil {
		t.Fatal("Expected markup, got nil")
	}

	if len(result.Markup.InlineKeyboard) != 1 || len(result.Markup.InlineKeyboard[0]) != 1 {
		t.Fatal("Expected one button in markup")
	}

	backButton := result.Markup.InlineKeyboard[0][0]
	if backButton.CallbackData == nil || *backButton.CallbackData != CallbackSetup {
		t.Errorf("Expected back button with CallbackSetup, got %v", backButton.CallbackData)
	}

	// Verify captured selection
	if capturedSelection == nil {
		t.Fatal("Expected user selection to be updated")
	}

	if !capturedSelection.CustomText {
		t.Error("Expected CustomText to be true")
	}

	if capturedSelection.ReminderMessage != "NLP_MODE" {
		t.Errorf("Expected ReminderMessage to be 'NLP_MODE', got '%s'", capturedSelection.ReminderMessage)
	}
}

func TestHandleNlpTextProcessing_Success(t *testing.T) {
	// Create a test user
	user := &entities.User{
		ID:       123,
		Language: LangEN,
	}

	// Create expected selection result
	expectedSelection := entities.NewUserSelection()
	expectedSelection.ReminderMessage = "Test reminder"

	// Mock NLP service
	mockNLP := &mockNLPService{
		shouldFail: false,
		result:     expectedSelection,
	}

	// Mock functions
	var capturedUserID int64
	var capturedSelection *entities.UserSelection
	mockCreateReminder := func(userID int64, selection *entities.UserSelection) (*entities.Reminder, error) {
		capturedUserID = userID
		capturedSelection = selection
		return &entities.Reminder{}, nil
	}

	var clearedUserID int64
	mockClearUserSelection := func(userID int64) error {
		clearedUserID = userID
		return nil
	}

	// Call the function
	result, err := HandleNlpTextProcessing(
		"Test reminder text",
		user,
		mockNLP,
		mockCreateReminder,
		mockClearUserSelection,
	)

	// Verify no error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Verify reminder was created with correct parameters
	if capturedUserID != user.ID {
		t.Errorf("Expected userID %d, got %d", user.ID, capturedUserID)
	}

	if capturedSelection != expectedSelection {
		t.Error("Expected captured selection to match expected selection")
	}

	// Verify user selection was cleared
	if clearedUserID != user.ID {
		t.Errorf("Expected cleared userID %d, got %d", user.ID, clearedUserID)
	}
}

func TestHandleNlpTextProcessing_NLPFailure(t *testing.T) {
	// Create a test user
	user := &entities.User{
		ID:       123,
		Language: LangEN,
	}

	// Mock NLP service that fails
	mockNLP := &mockNLPService{
		shouldFail: true,
	}

	// Mock functions (should not be called)
	mockCreateReminder := func(userID int64, selection *entities.UserSelection) (*entities.Reminder, error) {
		t.Error("CreateReminder should not be called when NLP fails")
		return nil, nil
	}

	mockClearUserSelection := func(userID int64) error {
		t.Error("ClearUserSelection should not be called when NLP fails")
		return nil
	}

	// Call the function
	result, err := HandleNlpTextProcessing(
		"Test reminder text",
		user,
		mockNLP,
		mockCreateReminder,
		mockClearUserSelection,
	)

	// Verify no error (errors are handled gracefully)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify result contains error message
	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	s := T(user.Language)
	expectedText := "❌ " + s.MsgParsingFailed + "\n\n" + s.NlpEnterText
	if result.Text != expectedText {
		t.Errorf("Expected error message, got: %s", result.Text)
	}

	// Verify markup has back button
	if result.Markup == nil {
		t.Fatal("Expected markup, got nil")
	}

	if len(result.Markup.InlineKeyboard) != 1 || len(result.Markup.InlineKeyboard[0]) != 1 {
		t.Fatal("Expected one button in markup")
	}
}
