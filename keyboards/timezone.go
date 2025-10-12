package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func HandleTimezoneSelection(user *entities.User, url string) (*SelectionResult, error) {
	s := T(user.Language)

	btn := tgbotapi.NewInlineKeyboardButtonURL(s.MsgTimezoneAutoDetect, url)
	backBtn := tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu)
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn),
		tgbotapi.NewInlineKeyboardRow(backBtn),
	)
	return &SelectionResult{Text: s.MsgTimezoneAutoDetectDescr, Markup: &markup}, nil
}
