package keyboards

import (
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
