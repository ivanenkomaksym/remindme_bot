package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
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
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Once), entities.Once.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Daily), entities.Daily.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Weekly), entities.Weekly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Monthly), entities.Monthly.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Interval), entities.Interval.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.SpacedBasedRepetition), entities.SpacedBasedRepetition.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMyReminders, CallbackRemindersList),
		),
	)

	return &mainMenu
}
