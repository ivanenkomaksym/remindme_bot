package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

// Premium menu callback data constants
const (
	CallbackPremiumUpgrade = "premium_upgrade"
)

// HandlePremiumSelection handles premium usage menu selections
func HandlePremiumSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, userUsage *entities.PremiumUsage) (*SelectionResult, error) {
	s := T(userEntity.Language)

	switch callbackData {
	case CallbackAccountViewPremium:
		// Show premium usage information
		var text string
		if userUsage != nil {
			text = FormatPremiumUsageInfo(userUsage, userEntity.Language)
		} else {
			text = s.PremiumTitle + "\n\n" + s.PremiumLoadError
		}

		// Create markup with upgrade button and back button
		upgradeBtn := tgbotapi.NewInlineKeyboardButtonData(s.PremiumUpgradeBtn, CallbackPremiumUpgrade)
		backBtn := tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu)

		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(upgradeBtn),
			tgbotapi.NewInlineKeyboardRow(backBtn),
		)

		return &SelectionResult{
			Text:   text,
			Markup: &markup,
		}, nil

	case CallbackPremiumUpgrade:
		// Show upgrade coming soon message
		return &SelectionResult{
			Text:   s.PremiumUpgradeComingSoon,
			Markup: GetAccountMenuMarkup(userEntity.Language),
		}, nil

	default:
		// Return to account menu for unhandled premium callbacks
		return &SelectionResult{
			Text:   s.MsgParsingFailed,
			Markup: GetAccountMenuMarkup(userEntity.Language),
		}, nil
	}
}

// IsPremiumCallback checks if the callback data is for premium management
func IsPremiumCallback(callbackData string) bool {
	return callbackData == CallbackAccountViewPremium ||
		callbackData == CallbackPremiumUpgrade
}
