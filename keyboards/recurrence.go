package keyboards

import (
	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
)

func HandleRecurrenceTypeSelection(callbackData string,
	user *entities.User,
	userSelection *entities.UserSelection) (*SelectionResult, error) {
	recurrenceType, err := entities.ToRecurrenceType(callbackData)
	if err != nil {
		return nil, err
	}

	userSelection.RecurrenceType = recurrenceType
	userSelection.IsWeekly = (recurrenceType == entities.Weekly)

	s := T(user.Language)
	switch recurrenceType {
	case entities.Once:
		return &SelectionResult{Text: s.MsgSelectDate, Markup: nil}, nil
	case entities.Daily:
		return &SelectionResult{Text: s.MsgSelectTime, Markup: GetHourRangeMarkup(user.Language)}, nil
	case entities.Weekly:
		return &SelectionResult{Text: s.MsgSelectWeekdays, Markup: GetWeekRangeMarkup(userSelection.WeekOptions, user.Language)}, nil
	case entities.Monthly:
		return &SelectionResult{Text: s.MsgSelectTime, Markup: GetHourRangeMarkup(user.Language)}, nil
	case entities.Interval:
		return &SelectionResult{Text: s.MsgSelectTime, Markup: GetHourRangeMarkup(user.Language)}, nil
	case entities.Custom:
		return &SelectionResult{Text: s.MsgEnterCustomTime, Markup: GetHourRangeMarkup(user.Language)}, nil
	}

	return nil, nil
}
