package keyboards

import (
	"fmt"
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

func GetWeekRangeMarkup(currentOptions [7]bool) *tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboard []tgbotapi.InlineKeyboardButton

	for idx, day := range LongDayNames {
		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardButtonData(buttonText(day, currentOptions[idx]), day))
	}
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardButtonData("Select", CallbackWeekSelect))

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{inlineKeyboard},
	}
}

func HandleWeekSelection(callbackData string, msg *tgbotapi.EditMessageTextConfig, currentOptions *[7]bool) *tgbotapi.InlineKeyboardMarkup {
	for idx, day := range LongDayNames {
		if callbackData == day {
			currentOptions[idx] = !currentOptions[idx]
		}
	}
	msg.Text = "Select weekdays"

	if callbackData == CallbackWeekSelect {
		msg.Text = fmt.Sprintf("Selected options: %v", *currentOptions)
		return nil
	}

	return GetWeekRangeMarkup(*currentOptions)
}

func buttonText(text string, opt bool) string {
	if opt {
		return "✅ " + text
	}

	return text
}
