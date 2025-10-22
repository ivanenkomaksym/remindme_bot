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
		// Featured: AI-powered text input (full width for prominence)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnNlpTextInput, CallbackNlpTextInput),
		),
		// Quick reminders: Once and Daily (most common options)
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Once), entities.Once.String()),
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Daily), entities.Daily.String()),
		),
		// Recurring reminders: Weekly and Monthly
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Weekly), entities.Weekly.String()),
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Monthly), entities.Monthly.String()),
		),
		// Advanced reminders: Interval and Spaced repetition
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.Interval), entities.Interval.String()),
			tgbotapi.NewInlineKeyboardButtonData(RecurrenceTypeLabel(lang, entities.SpacedBasedRepetition), entities.SpacedBasedRepetition.String()),
		),
		// Management and navigation
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(s.BtnMyReminders, CallbackRemindersList),
			tgbotapi.NewInlineKeyboardButtonData(s.BtnBack, CallbackBackToMainMenu),
		),
	)

	return &setupMenu
}
