package keyboards

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

// The callback data prefixes help to parse the user's message selection.
const (
	CallbackPrefixMessage  = "msg_"
	CallbackMessageCustom  = "msg_custom"
	CallbackMessageConfirm = "msg_confirm"
)

var DefaultMessages = []string{
	"Time to take a break!",
	"Don't forget your medication",
	"Check your email",
	"Drink some water",
	"Stand up and stretch",
	"Review your tasks",
}

func IsMessageSelectionCallback(callbackData string) bool {
	return strings.HasPrefix(callbackData, CallbackPrefixMessage)
}

func GetMessageSelectionMarkup() *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Add default message options
	for i, msg := range DefaultMessages {
		callbackData := CallbackPrefixMessage + string(rune(i))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(msg, callbackData)))
	}

	// Add custom message and confirm options
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✏️ Custom Message", CallbackMessageCustom),
		tgbotapi.NewInlineKeyboardButtonData("✅ Confirm", CallbackMessageConfirm),
	))

	// Add back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("← Back", "back_to_time"),
	))

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func HandleMessageSelection(callbackData string, msg *tgbotapi.EditMessageTextConfig, userState *types.UserSelectionState) *tgbotapi.InlineKeyboardMarkup {
	if strings.HasPrefix(callbackData, CallbackPrefixMessage) {
		// Extract message index
		msgIndex := int(callbackData[len(CallbackPrefixMessage)])
		if msgIndex >= 0 && msgIndex < len(DefaultMessages) {
			userState.ReminderMessage = DefaultMessages[msgIndex]
		}
		msg.Text = "Select your reminder message:"
		return GetMessageSelectionMarkup()
	}

	if callbackData == CallbackMessageCustom {
		msg.Text = "Please type your custom reminder message:"
		return nil
	}

	if callbackData == CallbackMessageConfirm {
		msg.Text = FormatReminderConfirmation(userState)
		return nil
	}

	return GetMessageSelectionMarkup()
}

func FormatReminderConfirmation(userState *types.UserSelectionState) string {
	confirmation := "✅ Reminder Set!\n\n"
	confirmation += "📅 Frequency: " + userState.RecurrenceType.String() + "\n"

	if userState.IsWeekly {
		confirmation += "📆 Days: "
		days := []string{}
		for i, selected := range userState.WeekOptions {
			if selected {
				days = append(days, LongDayNames[i])
			}
		}
		if len(days) > 0 {
			confirmation += strings.Join(days, ", ")
		} else {
			confirmation += "None selected"
		}
		confirmation += "\n"
	}

	confirmation += "⏰ Time: " + userState.SelectedTime + "\n"
	confirmation += "💬 Message: " + userState.ReminderMessage + "\n\n"
	confirmation += "Your reminder has been scheduled!"

	return confirmation
}
