package keyboards

import (
	"slices"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var LongDayNames = []string{time.Monday.String(),
	time.Tuesday.String(),
	time.Wednesday.String(),
	time.Thursday.String(),
	time.Friday.String(),
	time.Saturday.String(),
	time.Sunday.String()}

// The callback data prefixes help to parse the user's selection.
const (
	CallbackWeekSelect = "week_select"
)

func IsWeekSelectionCallback(callbackData string) bool {
	if slices.Contains(LongDayNames, callbackData) {
		return true
	}

	if callbackData == CallbackWeekSelect {
		return true
	}

	return false
}

func GetWeekRangeMarkup(currentOptions [7]bool) tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboard []tgbotapi.InlineKeyboardButton

	for _, day := range LongDayNames {
		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardButtonData(buttonText(day, currentOptions[0]), day))
	}
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardButtonData("Select", CallbackWeekSelect))

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{inlineKeyboard},
	}
}

func ProcessWeekSelection(callbackQueryData string) {

}

func buttonText(text string, opt bool) string {
	if opt {
		return "✅ " + text
	}

	return "❌ " + text
}
