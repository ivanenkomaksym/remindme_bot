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

	s := T(userState.Language)
	switch recurrenceType {
	case models.Daily:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(userState.Language), nil
	case models.Weekly:
		msg.Text = s.MsgSelectWeekdays
		return GetWeekRangeMarkup(userState.WeekOptions, userState.Language), nil
	case models.Monthly:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(userState.Language), nil
	case models.Interval:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(userState.Language), nil
	case models.Custom:
		msg.Text = s.MsgEnterCustomTime
		return GetHourRangeMarkup(userState.Language), nil
	}

	return nil, nil
}
