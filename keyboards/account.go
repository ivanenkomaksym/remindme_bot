package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func HandleAccountSelection(user *tgbotapi.User, callbackData string, userEntity *entities.User, timezoneURL string) (*SelectionResult, error) {
	switch callbackData {
	case CallbackAccountChangeLanguage:
		// Show language selection menu
		return &SelectionResult{
			Text:   "Select language / Оберіть мову:",
			Markup: GetLanguageSelectionMarkup(),
		}, nil
	case CallbackAccountChangeTimezone:
		// Show timezone setup
		return HandleTimezoneSelection(userEntity, timezoneURL)
	case CallbackAccountBackToMenu:
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
