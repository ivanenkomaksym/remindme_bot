package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The callback data prefixes help to parse the user's selection.
const (
	CallbackWeekSelect = "week_select"
	CallbackWeekDay    = "week_day:"
)

func IsWeekSelectionCallback(callbackData string) bool {
	if callbackData == CallbackWeekSelect {
		return true
	}
	if stringsHasPrefix(callbackData, CallbackWeekDay) {
		return true
	}
	return false
}

func GetWeekRangeMarkup(currentOptions [7]bool, lang string) *tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboard [][]tgbotapi.InlineKeyboardButton
	s := T(lang)
	weekdayNames := s.WeekdayNames

	for idx, day := range weekdayNames {
		callback := fmt.Sprintf("%s%d", CallbackWeekDay, idx)
		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText(day, currentOptions[idx]), callback)))
	}
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnSelect, CallbackWeekSelect)))
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, MainMenu)))

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
}

func HandleWeekSelection(callbackData string, msg *tgbotapi.EditMessageTextConfig, currentOptions *[7]bool, lang string) *tgbotapi.InlineKeyboardMarkup {
	if stringsHasPrefix(callbackData, CallbackWeekDay) {
		var idx int
		_, _ = fmt.Sscanf(callbackData[len(CallbackWeekDay):], "%d", &idx)
		if idx >= 0 && idx < 7 {
			currentOptions[idx] = !currentOptions[idx]
		}
	}
	s := T(lang)
	msg.Text = s.MsgSelectWeekdays

	if callbackData == CallbackWeekSelect {
		msg.Text = s.MsgSelectTimeWeekly
		return GetHourRangeMarkup(lang)
	}

	return GetWeekRangeMarkup(*currentOptions, lang)
}

func buttonText(text string, opt bool) string {
	if opt {
		return "✅ " + text
	}

	return text
}

// tiny helper to avoid importing strings; keep consistent spacing and style
func stringsHasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
