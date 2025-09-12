package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
)

const (
	MainMenu = "main_menu"
)

func IsMainMenuSelection(callbackData string) bool {
	return callbackData == MainMenu
}

func GetMainMenuMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	s := T(lang)
	mainMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, models.Daily), models.Daily.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, models.Weekly), models.Weekly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, models.Monthly), models.Monthly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, models.Interval), models.Interval.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, models.Custom), models.Custom.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMyReminders, CallbackRemindersList),
		),
	)

	return &mainMenu
}
