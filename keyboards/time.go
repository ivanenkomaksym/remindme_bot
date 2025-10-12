package keyboards

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
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
	user *entities.User,
	userSelection *entities.UserSelection) *SelectionResult {
	s := T(user.Language)
	switch {
	case strings.Contains(callbackData, CallbackTimeStart):
		return &SelectionResult{Text: s.MsgSelectTime, Markup: GetHourRangeMarkup(user.Language)}

	case strings.Contains(callbackData, CallbackPrefixHourRange):
		startHour := 0
		fmt.Sscanf(callbackData[len(CallbackPrefixHourRange):], "%d", &startHour)
		return &SelectionResult{Text: fmt.Sprintf(s.MsgSelectWithinHour, startHour, (startHour+4)%24), Markup: GetMinuteRangeMarkup(startHour, user.Language)}

	case strings.Contains(callbackData, CallbackPrefixMinuteRange):
		startHour := 0
		fmt.Sscanf(callbackData[len(CallbackPrefixMinuteRange):], "%d", &startHour)
		return &SelectionResult{Text: fmt.Sprintf(s.MsgSelectWithinHour, startHour, (startHour+1)%24), Markup: GetSpecificTimeMarkup(startHour, user.Language)}

	case strings.Contains(callbackData, CallbackPrefixSpecificTime):
		timeStr := callbackData[len(CallbackPrefixSpecificTime):]
		userSelection.SelectedTime = timeStr
		return &SelectionResult{Text: s.MsgSelectMessage, Markup: GetMessageSelectionMarkup(user.Language)}

	case strings.Contains(callbackData, CallbackPrefixCustom):
		userSelection.CustomTime = true
		return &SelectionResult{Text: s.MsgEnterCustomTime, Markup: nil}
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
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, SetupMenu),
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

func HandleCustomTimeSelection(text string,
	msg *tgbotapi.MessageConfig,
	user *entities.User,
	userSelection *entities.UserSelection) *SelectionResult {
	_, _, ok := scheduler.ParseHourMinute(text)

	outputText := ""
	markup := &tgbotapi.InlineKeyboardMarkup{}

	if !ok {
		s := T(user.Language)
		outputText = fmt.Sprintf("%s. %s", s.MsgInvalidTimeFormat, s.MsgEnterCustomTime)
		markup = GetHourRangeMarkup(user.Language)
	} else {
		userSelection.SelectedTime = text
		s := T(user.Language)
		outputText = s.MsgSelectMessage
		markup = GetMessageSelectionMarkup(user.Language)
	}

	return &SelectionResult{Text: outputText, Markup: markup}
}
