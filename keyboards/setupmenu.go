package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

const (
	SetupMenu            = "setup_menu"
	CallbackNlpTextInput = "nlp_text_input"
)

func IsSetupMenuSelection(callbackData string) bool {
	return callbackData == SetupMenu
}

func IsNlpTextInputCallback(callbackData string) bool {
	return callbackData == CallbackNlpTextInput
}

func GetSetupMenuMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	s := T(lang)
	setupMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnNlpTextInput, CallbackNlpTextInput),
		),
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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu),
		),
	)

	return &setupMenu
}
