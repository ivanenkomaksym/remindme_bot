package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// Navigation menu callback data constants
const (
	CallbackList    = "nav_list"
	CallbackSetup   = "nav_setup"
	CallbackAccount = "nav_account"
	// Account management callbacks
	CallbackAccountChangeLanguage = "acc_change_lang"
	CallbackAccountChangeTimezone = "acc_change_tz"
	// General back to main menu callback
	CallbackBackToMainMenu = "back_to_main"
	// Timezone selection callbacks
	CallbackTimezoneManual = "tz_manual"
	CallbackTimezoneSelect = "tz_select_"
)

func IsMainMenuSelection(callbackData string) bool {
	return callbackData == CallbackBackToMainMenu
}

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
			tgbotapi.NewInlineKeyboardButtonData("💼 "+s.NavAccount, CallbackAccount),
		),
	)

	return &navMenu
} // IsNavigationCallback checks if the callback data is for navigation
func IsNavigationCallback(callbackData string) bool {
	return callbackData == CallbackList ||
		callbackData == CallbackSetup ||
		callbackData == CallbackAccount
}

// GetAccountMenuMarkup returns the account management keyboard
func GetAccountMenuMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	s := T(lang)

	accountMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.AccChangeLanguage, CallbackAccountChangeLanguage),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.AccChangeTimezone, CallbackAccountChangeTimezone),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu),
		),
	)

	return &accountMenu
}

// FormatAccountInfo formats user account information for display
func FormatAccountInfo(user *entities.User, lang string) string {
	s := T(lang)

	username := user.UserName
	if username == "" {
		username = s.AccNoUsername
	}

	language := "English"
	if user.Language == LangUK {
		language = "Українська"
	}

	timezone := s.AccNoTimezone
	if user.GetLocation() != nil {
		timezone = user.GetLocation().String()
	}

	createdAt := user.CreatedAt.Format("2006-01-02 15:04")

	return fmt.Sprintf("%s\n\n"+
		"📝 %s: @%s\n"+
		"🌐 %s: %s\n"+
		"🕐 %s: %s\n"+
		"📅 %s: %s",
		s.AccTitle,
		s.AccUsername, username,
		s.AccLanguage, language,
		s.AccTimezone, timezone,
		s.AccCreatedAt, createdAt)
}

// IsAccountCallback checks if the callback data is for account management
func IsAccountCallback(callbackData string) bool {
	return callbackData == CallbackAccountChangeLanguage ||
		callbackData == CallbackAccountChangeTimezone
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
