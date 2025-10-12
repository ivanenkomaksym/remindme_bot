package keyboards

import (
	"fmt"
	"slices"
	"time"

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

func GetWeekRangeMarkup(weekdays []time.Weekday, lang string) *tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboard [][]tgbotapi.InlineKeyboardButton
	s := T(lang)
	weekdayNames := s.WeekdayNames

	keys := T("en").WeekdayNamesShort

	for i := 1; i < 8; i++ {
		weekday := time.Weekday(i % 7)
		name := weekdayNames[weekday]
		key := keys[weekday]
		callback := fmt.Sprintf("%s%s", CallbackWeekDay, key)
		inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText(name, slices.Contains(weekdays, weekday)), callback)))
	}
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnSelect, CallbackWeekSelect)))
	inlineKeyboard = append(inlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, SetupMenu)))

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
}

func HandleWeekSelection(callbackData string, weekdays *[]time.Weekday, lang string) *SelectionResult {
	if stringsHasPrefix(callbackData, CallbackWeekDay) {
		var weekdayStr string
		_, _ = fmt.Sscanf(callbackData[len(CallbackWeekDay):], "%s", &weekdayStr)
		weekday, ok := WeekdayNameToKeyMap[weekdayStr]
		if !ok {
			return nil
		}

		if slices.Contains(*weekdays, weekday) {
			// Remove
			*weekdays = slices.Delete(*weekdays, slices.Index(*weekdays, weekday), slices.Index(*weekdays, weekday)+1)
		} else {
			// Add
			*weekdays = append(*weekdays, weekday)
			slices.Sort(*weekdays)
		}
	}
	s := T(lang)
	if callbackData == CallbackWeekSelect {
		return &SelectionResult{Text: s.MsgSelectTimeWeekly, Markup: GetHourRangeMarkup(lang)}
	}
	return &SelectionResult{Text: s.MsgSelectWeekdays, Markup: GetWeekRangeMarkup(*weekdays, lang)}
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
