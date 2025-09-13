package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanenkomaksym/remindme_bot/models"
)

func HandleRecurrenceTypeSelection(callbackData string,
	msg *tgbotapi.EditMessageTextConfig,
	user *models.User,
	userSelection *models.UserSelection) (*tgbotapi.InlineKeyboardMarkup, error) {
	recurrenceType, err := models.ToRecurrenceType(callbackData)
	if err != nil {
		return nil, err
	}

	userSelection.RecurrenceType = recurrenceType
	userSelection.IsWeekly = (recurrenceType == models.Weekly)

	s := T(user.Language)
	switch recurrenceType {
	case models.Daily:
		msg.Text = s.MsgSelectTime
		return GetHourRangeMarkup(user.Language), nil
	case models.Weekly:
		msg.Text = s.MsgSelectWeekdays
		return GetWeekRangeMarkup(userSelection.WeekOptions, user.Language), nil
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
