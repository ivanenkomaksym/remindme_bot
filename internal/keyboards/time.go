package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The callback data prefixes help to parse the user's selection.
const (
	// Represents the start of a 4-hour range selection.
	CallbackPrefixHourRange = "time_hour_range:"
	// Represents the start of a 1-hour range selection.
	CallbackPrefixMinuteRange = "time_minute_range:"
	// Represents a specific 15-minute time selection.
	CallbackPrefixSpecificTime = "time_specific:"
	// Represents the custom time selection option.
	CallbackPrefixCustom = "time_custom:"
)

// getHourRangeMarkup generates the first level of the menu (4-hour ranges).
func GetHourRangeMarkup() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var buttons []tgbotapi.InlineKeyboardButton

	// Create buttons for 4-hour blocks: 0:00-4:00, 4:00-8:00, etc.
	for i := 0; i < 24; i += 4 {
		start := fmt.Sprintf("%02d:00", i)
		end := fmt.Sprintf("%02d:00", (i+4)%24)
		text := fmt.Sprintf("%s-%s", start, end)
		callbackData := fmt.Sprintf("%s%d", CallbackPrefixHourRange, i)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(text, callbackData))
	}

	rows = append(rows, buttons)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Custom", CallbackPrefixCustom),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// getMinuteRangeMarkup generates the second level of the menu (1-hour ranges).
func GetMinuteRangeMarkup(startHour int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var buttons []tgbotapi.InlineKeyboardButton

	// The loop should create 4 buttons for the next 4 hours
	for i := 0; i < 4; i++ {
		currentHour := (startHour + i) % 24
		nextHour := (currentHour + 1) % 24
		text := fmt.Sprintf("%02d:00-%02d:00", currentHour, nextHour)
		callbackData := fmt.Sprintf("%s%d", CallbackPrefixMinuteRange, currentHour)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(text, callbackData))
	}

	rows = append(rows, buttons)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// getSpecificTimeMarkup generates the third and final level of the menu (15-minute intervals).
func GetSpecificTimeMarkup(startHour int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var buttons []tgbotapi.InlineKeyboardButton

	// Create buttons for 15-minute intervals within the selected hour.
	for i := 0; i < 60; i += 15 {
		text := fmt.Sprintf("%02d:%02d", startHour, i)
		callbackData := fmt.Sprintf("%s%02d:%02d", CallbackPrefixSpecificTime, startHour, i)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(text, callbackData))
	}

	// Add the next full hour as a final option.
	nextHourText := fmt.Sprintf("%02d:00", (startHour+1)%24)
	nextHourCallbackData := fmt.Sprintf("%s%02d:00", CallbackPrefixSpecificTime, (startHour+1)%24)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(nextHourText, nextHourCallbackData))

	rows = append(rows, buttons)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Custom", CallbackPrefixCustom),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
