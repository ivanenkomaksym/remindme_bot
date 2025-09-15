package keyboards

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	LangEN = "en"
	LangUK = "uk"

	CallbackLangPrefix = "lang:"
)

func GetLanguageSelectionMarkup() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("English", CallbackLangPrefix+LangEN),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Українська", CallbackLangPrefix+LangUK),
		),
	)
	return &menu
}

func IsLanguageSelectionCallback(data string) bool {
	return len(data) >= len(CallbackLangPrefix) && data[:len(CallbackLangPrefix)] == CallbackLangPrefix
}

func ParseLanguageFromCallback(data string) string {
	if !IsLanguageSelectionCallback(data) {
		return ""
	}
	return data[len(CallbackLangPrefix):]
}

// MapTelegramLanguageCodeToSupported maps Telegram's User.LanguageCode (e.g. "en", "en-US", "uk-UA")
// to a supported language code. Returns (code, true) if supported; otherwise ("", false).
func MapTelegramLanguageCodeToSupported(code string) (string, bool) {
	if code == "" {
		return "", false
	}
	lower := strings.ToLower(code)
	if strings.HasPrefix(lower, "uk") {
		return LangUK, true
	}
	if strings.HasPrefix(lower, "en") {
		return LangEN, true
	}
	return "", false
}
