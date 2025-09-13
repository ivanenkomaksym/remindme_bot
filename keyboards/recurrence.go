package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/models"
	"github.com/ivanenkomaksym/remindme_bot/types"
)

func HandleRecurrenceTypeSelection(callbackData string,
	msg *tgbotapi.EditMessageTextConfig,
	user *models.User,
	userState *types.UserSelectionState) (*tgbotapi.InlineKeyboardMarkup, error) {
	recurrenceType, err := models.ToRecurrenceType(callbackData)
	if err != nil {
		return nil, err
	}

	userState.RecurrenceType = recurrenceType
	userState.IsWeekly = (recurrenceType == models.Weekly)

	s := T(user.Language)
	switch recurrenceType {
	case models.Daily:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(user.Language), nil
	case models.Weekly:
		msg.Text = s.MsgSelectWeekdays
		return GetWeekRangeMarkup(userState.WeekOptions, user.Language), nil
	case models.Monthly:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(user.Language), nil
	case models.Interval:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(user.Language), nil
	case models.Custom:
		msg.Text = s.MsgEnterCustomTime
		return GetHourRangeMarkup(user.Language), nil
	}

	return nil, nil
}
