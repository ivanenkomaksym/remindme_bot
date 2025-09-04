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

func GetMainMenuMarkup() *tgbotapi.InlineKeyboardMarkup {
	mainMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Daily.String(), models.Daily.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Weekly.String(), models.Weekly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Monthly.String(), models.Monthly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Interval.String(), models.Interval.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(models.Custom.String(), models.Custom.String()),
		),
	)

	return &mainMenu
}
