package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func HandleAccountSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, timezoneURL string, userUsage *entities.PremiumUsage) (*SelectionResult, error) {
	switch callbackData {
	case CallbackAccountChangeLanguage:
		// Show language selection menu
		s := T(userEntity.Language)
		return &SelectionResult{
			Text:   s.LanguageSelectPrompt,
			Markup: GetLanguageSelectionMarkup(userEntity.Language),
		}, nil
	case CallbackAccountChangeTimezone:
		// Show timezone setup
		return HandleTimezoneSelection(userEntity, timezoneURL)
	case CallbackAccountViewPremium:
		// Show premium usage information
		s := T(userEntity.Language)
		var text string
		if userUsage != nil {
			text = FormatPremiumUsageInfo(userUsage, userEntity.Language)
		} else {
			text = s.PremiumTitle + "\n\n" + s.PremiumLoadError
		}
		return &SelectionResult{
			Text:   text,
			Markup: GetAccountMenuMarkup(userEntity.Language),
		}, nil
	case CallbackBackToMainMenu:
		// Return to main navigation menu
		s := T(userEntity.Language)
		return &SelectionResult{
			Text:   s.Welcome + "\n\n" + s.NavChooseOption,
			Markup: GetNavigationMenuMarkup(userEntity.Language),
		}, nil
	default:
		s := T(userEntity.Language)
		return &SelectionResult{
			Text:   s.MsgParsingFailed,
			Markup: GetAccountMenuMarkup(userEntity.Language),
		}, nil
	}
}
