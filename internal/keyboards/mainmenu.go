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
			tgbotapi.NewInlineKeyboardButtonData(s.BtnDaily, models.Daily.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnWeekly, models.Weekly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMonthly, models.Monthly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnInterval, models.Interval.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnCustom, models.Custom.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMyReminders, CallbackRemindersList),
		),
	)

	return &mainMenu
}
