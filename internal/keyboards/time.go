package keyboards

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

// The callback data prefixes help to parse the user's selection.
const (
	CallbackTimeStart = "time_start"
	// Represents the start of a 4-hour range selection.
	CallbackPrefixHourRange = "time_hour_range:"
	// Represents the start of a 1-hour range selection.
	CallbackPrefixMinuteRange = "time_minute_range:"
	// Represents a specific 15-minute time selection.
	CallbackPrefixSpecificTime = "time_specific:"
	// Represents the custom time selection option.
	CallbackPrefixCustom = "time_custom:"
)

func IsTimeSelectionCallback(callbackData string) bool {
	return strings.HasPrefix(callbackData, "time_")
}

func HandleTimeSelection(callbackData string,
	msg *tgbotapi.EditMessageTextConfig,
	userState *types.UserSelectionState) *tgbotapi.InlineKeyboardMarkup {
	s := T(userState.Language)
	switch {
	case strings.Contains(callbackData, CallbackTimeStart):
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(userState.Language)

	case strings.Contains(callbackData, CallbackPrefixHourRange):
		// User selected a 4-hour range, show 1-hour ranges
		startHour := 0
		fmt.Sscanf(callbackData[len(CallbackPrefixHourRange):], "%d", &startHour)
		msg.Text = fmt.Sprintf(s.MsgSelectWithinHour, startHour, (startHour+4)%24)
		return GetMinuteRangeMarkup(startHour, userState.Language)

	case strings.Contains(callbackData, CallbackPrefixMinuteRange):
		// User selected a 1-hour range, show 15-minute intervals
		startHour := 0
		fmt.Sscanf(callbackData[len(CallbackPrefixMinuteRange):], "%d", &startHour)
		msg.Text = fmt.Sprintf(s.MsgSelectWithinHour, startHour, (startHour+1)%24)
		return GetSpecificTimeMarkup(startHour, userState.Language)

	case strings.Contains(callbackData, CallbackPrefixSpecificTime):
		// User selected a specific time, go to message selection
		timeStr := callbackData[len(CallbackPrefixSpecificTime):]
		userState.SelectedTime = timeStr
		msg.Text = s.MsgSelectMessage
		return GetMessageSelectionMarkup(userState.Language)

	case strings.Contains(callbackData, CallbackPrefixCustom):
		// User wants custom time input
		userState.CustomTime = true
		msg.Text = s.MsgEnterCustomTime
		return nil
	}

	return nil
}

// getHourRangeMarkup generates the first level of the menu (4-hour ranges).
func GetHourRangeMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	s := T(lang)
	// Create buttons for 4-hour blocks: 0:00-4:00, 4:00-8:00, etc.
	for i := 0; i < 24; i += 4 {
		start := fmt.Sprintf("%02d:00", i)
		end := fmt.Sprintf("%02d:00", (i+4)%24)
		text := fmt.Sprintf("%s-%s", start, end)
		callbackData := fmt.Sprintf("%s%d", CallbackPrefixHourRange, i)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, callbackData)))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnCustomTime, CallbackPrefixCustom),
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, MainMenu),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

// getMinuteRangeMarkup generates the second level of the menu (1-hour ranges).
func GetMinuteRangeMarkup(startHour int, lang string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	s := T(lang)
	// The loop should create 4 buttons for the next 4 hours
	for i := 0; i < 4; i++ {
		currentHour := (startHour + i) % 24
		nextHour := (currentHour + 1) % 24
		text := fmt.Sprintf("%02d:00-%02d:00", currentHour, nextHour)
		callbackData := fmt.Sprintf("%s%d", CallbackPrefixMinuteRange, currentHour)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, callbackData)))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnCustomTime, CallbackPrefixCustom),
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackTimeStart),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

// getSpecificTimeMarkup generates the third and final level of the menu (15-minute intervals).
func GetSpecificTimeMarkup(startHour int, lang string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var buttons []tgbotapi.InlineKeyboardButton
	s := T(lang)
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

	backData := fmt.Sprintf("%s%d", CallbackPrefixHourRange, 4*(startHour/4))

	rows = append(rows, buttons)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnCustomTime, CallbackPrefixCustom),
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, backData),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func HadleCustomTimeSelection(text string,
	msg *tgbotapi.MessageConfig,
	userState *types.UserSelectionState) *tgbotapi.InlineKeyboardMarkup {
	if !isValidTimeFormat(text) {
		s := T(userState.Language)
		msg.Text = fmt.Sprintf("%s. %s", s.MsgInvalidTimeFormat, s.MsgEnterCustomTime)
		return GetHourRangeMarkup(userState.Language)
	} else {
		userState.SelectedTime = text
		s := T(userState.Language)
		msg.Text = s.MsgSelectMessage
		return GetMessageSelectionMarkup(userState.Language)
	}
}

// isValidTimeFormat checks if the input string is a valid time format (HH:MM)
func isValidTimeFormat(timeStr string) bool {
	// Accepts H:MM or HH:MM
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}
	hourStr, minStr := parts[0], parts[1]
	if len(hourStr) < 1 || len(hourStr) > 2 || len(minStr) != 2 {
		return false
	}
	hour := 0
	minute := 0
	var err error
	hour, err = strconv.Atoi(hourStr)
	if err != nil || hour < 0 || hour > 23 {
		return false
	}
	minute, err = strconv.Atoi(minStr)
	if err != nil || minute < 0 || minute > 59 {
		return false
	}
	return true
}
