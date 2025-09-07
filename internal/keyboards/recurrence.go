package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/internal/models"
	"github.com/ivanenkomaksym/remindme_bot/internal/types"
)

func HandleRecurrenceTypeSelection(callbackData string,
	msg *tgbotapi.EditMessageTextConfig,
	userState *types.UserSelectionState) (*tgbotapi.InlineKeyboardMarkup, error) {
	recurrenceType, err := models.ToRecurrenceType(callbackData)
	if err != nil {
		return nil, err
	}

	userState.RecurrenceType = recurrenceType
	userState.IsWeekly = (recurrenceType == models.Weekly)

	switch recurrenceType {
	case models.Daily:
		msg.Text = "Select time for daily reminders:"
		return GetHourRangeMarkup(), nil
	case models.Weekly:
		msg.Text = "Select time for weekly reminders:"
		return GetWeekRangeMarkup(userState.WeekOptions), nil
	case models.Monthly:
		msg.Text = "Select time for monthly reminders:"
		return GetHourRangeMarkup(), nil
	case models.Interval:
		msg.Text = "Select time for interval reminders:"
		return GetHourRangeMarkup(), nil
	case models.Custom:
		msg.Text = "Please type your custom time in HH:MM format (e.g., 14:05):"
		return GetHourRangeMarkup(), nil
	}

	return nil, nil
}
