package keyboards

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func HandleTimezoneSelection(user *entities.User, url string) (*SelectionResult, error) {
	s := T(user.Language)

	autoBtn := tgbotapi.NewInlineKeyboardButtonURL(s.MsgTimezoneAutoDetect, url)
	manualBtn := tgbotapi.NewInlineKeyboardButtonData(s.TzManualSelect, CallbackTimezoneManual)
	backBtn := tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(autoBtn),
		tgbotapi.NewInlineKeyboardRow(manualBtn),
		tgbotapi.NewInlineKeyboardRow(backBtn),
	)
	return &SelectionResult{Text: s.MsgTimezoneAutoDetectDescr, Markup: &markup}, nil
}

// GetTimezoneSelectionMarkup returns a keyboard with common timezone options
func GetTimezoneSelectionMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	s := T(lang)

	// Common timezones with their display names and IANA identifiers
	timezones := []struct {
		display string
		iana    string
	}{
		{"🇺🇸 New York (EST)", "America/New_York"},
		{"🇺🇸 Los Angeles (PST)", "America/Los_Angeles"},
		{"🇺🇸 Chicago (CST)", "America/Chicago"},
		{"🇬🇧 London (GMT)", "Europe/London"},
		{"🇫🇷 Paris (CET)", "Europe/Paris"},
		{"🇩🇪 Berlin (CET)", "Europe/Berlin"},
		{"🇺🇦 Kyiv (EET)", "Europe/Kiev"},
		{"🇷🇺 Moscow (MSK)", "Europe/Moscow"},
		{"🇯🇵 Tokyo (JST)", "Asia/Tokyo"},
		{"🇨🇳 Beijing (CST)", "Asia/Shanghai"},
		{"🇮🇳 Mumbai (IST)", "Asia/Kolkata"},
		{"🇦🇺 Sydney (AEDT)", "Australia/Sydney"},
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, tz := range timezones {
		btn := tgbotapi.NewInlineKeyboardButtonData(tz.display, CallbackTimezoneSelect+tz.iana)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	// Add back button
	backBtn := tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backBtn))

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// HandleManualTimezoneSelection handles the manual timezone selection flow
func HandleManualTimezoneSelection(callbackData string, lang string) (*SelectionResult, error) {
	s := T(lang)

	if callbackData == CallbackTimezoneManual {
		// Show timezone selection menu
		return &SelectionResult{
			Text:   s.TzSelectPrompt,
			Markup: GetTimezoneSelectionMarkup(lang),
		}, nil
	}

	if strings.HasPrefix(callbackData, CallbackTimezoneSelect) {
		// Extract timezone from callback data
		timezone := strings.TrimPrefix(callbackData, CallbackTimezoneSelect)

		// Return success message (the actual timezone setting would be handled by the use case)
		return &SelectionResult{
			Text:   s.MsgTimezoneSet + " " + timezone,
			Markup: GetNavigationMenuMarkup(lang),
		}, nil
	}

	return nil, nil
}

// IsTimezoneCallback checks if callback is for timezone selection
func IsTimezoneCallback(callbackData string) bool {
	return callbackData == CallbackTimezoneManual ||
		strings.HasPrefix(callbackData, CallbackTimezoneSelect)
}
