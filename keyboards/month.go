package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The callback data prefixes help to parse the user's selection for month days.
const (
	CallbackMonthSelect = "month_select"
	CallbackMonthDay    = "month_day:"
)

func IsMonthSelectionCallback(callbackData string) bool {
	if callbackData == CallbackMonthSelect {
		return true
	}
	if stringsHasPrefix(callbackData, CallbackMonthDay) {
		return true
	}
	return false
}

// GetMonthRangeMarkup renders a 4x7 grid for days 1..28 with multi-select support.
func GetMonthRangeMarkup(currentOptions [28]bool, lang string) *tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboard [][]tgbotapi.InlineKeyboardButton
	s := T(lang)

	// 4 rows of 7 days each
	day := 1
	for r := 0; r < 4; r++ {
		var row []tgbotapi.InlineKeyboardButton
		for c := 0; c < 7; c++ {
			label := fmt.Sprintf("%d", day)
			callback := fmt.Sprintf("%s%d", CallbackMonthDay, day)
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(buttonText(label, currentOptions[day-1]), callback))
			day++
		}
		inlineKeyboard = append(inlineKeyboard, row)
	}

	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnSelect, CallbackMonthSelect)))
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, MainMenu)))

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
}

// HandleMonthSelection toggles a selected day and returns the updated view or proceeds on select.
func HandleMonthSelection(callbackData string, currentOptions *[28]bool, lang string) *SelectionResult {
	if stringsHasPrefix(callbackData, CallbackMonthDay) {
		var day int
		_, _ = fmt.Sscanf(callbackData[len(CallbackMonthDay):], "%d", &day)
		if day >= 1 && day <= 28 {
			idx := day - 1
			currentOptions[idx] = !currentOptions[idx]
		}
	}
	s := T(lang)
	if callbackData == CallbackMonthSelect {
		return &SelectionResult{Text: s.MsgSelectTime, Markup: GetHourRangeMarkup(lang)}
	}
	return &SelectionResult{Text: s.MsgSelectDate, Markup: GetMonthRangeMarkup(*currentOptions, lang)}
}
