package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Navigation menu callback data constants
const (
	CallbackList    = "nav_list"
	CallbackSetup   = "nav_setup"
	CallbackAccount = "nav_account"
)

// GetNavigationMenuMarkup returns the main navigation menu keyboard
func GetNavigationMenuMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	s := T(lang)

	navMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 "+s.NavList, CallbackList),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ "+s.NavSetup, CallbackSetup),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 "+s.NavAccount, CallbackAccount),
		),
	)

	return &navMenu
} // IsNavigationCallback checks if the callback data is for navigation
func IsNavigationCallback(callbackData string) bool {
	return callbackData == CallbackList ||
		callbackData == CallbackSetup ||
		callbackData == CallbackAccount
}

// HandleNavigationCallback handles navigation menu callbacks
func HandleNavigationCallback(callbackData string) string {
	switch callbackData {
	case CallbackList:
		return "/list"
	case CallbackSetup:
		return "/setup"
	case CallbackAccount:
		return "/account"
	default:
		return ""
	}
}
