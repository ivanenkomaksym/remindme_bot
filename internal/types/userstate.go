package types

import "github.com/ivanenkomaksym/remindme_bot/internal/models"

// UserSelectionState holds the complete state of a user's reminder setup
type UserSelectionState struct {
	RecurrenceType  models.RecurrenceType
	IsWeekly        bool
	WeekOptions     [7]bool
	SelectedTime    string
	ReminderMessage string
	CustomText      bool
}
