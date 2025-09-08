package keyboards

import (
	"strings"
	"testing"
)

// TestGetHourRangeMarkup verifies the first level of the keyboard.
func TestGetHourRangeMarkup(t *testing.T) {
	markup := GetHourRangeMarkup()

	// Check the callback data of a specific button.
	button := markup.InlineKeyboard[2][0] // The "8:00-12:00" button.
	expectedCallback := "time_hour_range:8"
	if button.CallbackData != nil && *button.CallbackData != expectedCallback {
		t.Errorf("Expected callback data %s, got %s", expectedCallback, *button.CallbackData)
	}
}

// TestGetMinuteRangeMarkup verifies the second level of the keyboard.
func TestGetMinuteRangeMarkup(t *testing.T) {
	startHour := 8 // Simulating a user selecting "8:00-12:00"
	markup := GetMinuteRangeMarkup(startHour)

	// Check the text of a specific button.
	button := markup.InlineKeyboard[1][0] // The "9:00-10:00" button.
	expectedText := "09:00-10:00"
	if button.Text != expectedText {
		t.Errorf("Expected button text %s, got %s", expectedText, button.Text)
	}

	// Check the callback data of a specific button.
	expectedCallbackPrefix := "time_minute_range:"
	if !strings.HasPrefix(*button.CallbackData, expectedCallbackPrefix) {
		t.Errorf("Expected callback data to have prefix %s", expectedCallbackPrefix)
	}
}

// TestGetSpecificTimeMarkup verifies the third level of the keyboard.
func TestGetSpecificTimeMarkup(t *testing.T) {
	startHour := 9 // Simulating a user selecting "9:00-10:00"
	markup := GetSpecificTimeMarkup(startHour)

	// There should be two rows: one for times, one for controls.
	if len(markup.InlineKeyboard) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(markup.InlineKeyboard))
	}

	// The first row should have 5 buttons (9:00, 9:15, 9:30, 9:45, 10:00).
	if len(markup.InlineKeyboard[0]) != 5 {
		t.Errorf("Expected 5 buttons in the first row, got %d", len(markup.InlineKeyboard[0]))
	}

	// The last row should have two buttons ("Custom" and "← Back").
	if len(markup.InlineKeyboard[1]) != 2 {
		t.Errorf("Expected 2 buttons in the second row, got %d", len(markup.InlineKeyboard[1]))
	}

	// Check the text of a specific button.
	button := markup.InlineKeyboard[0][1] // The "9:15" button.
	expectedText := "09:15"
	if button.Text != expectedText {
		t.Errorf("Expected button text %s, got %s", expectedText, button.Text)
	}

	// Check the callback data of the last button ("10:00").
	lastButton := markup.InlineKeyboard[0][4]
	expectedCallback := "time_specific:10:00"
	if *lastButton.CallbackData != expectedCallback {
		t.Errorf("Expected callback data %s, got %s", expectedCallback, *lastButton.CallbackData)
	}
}
