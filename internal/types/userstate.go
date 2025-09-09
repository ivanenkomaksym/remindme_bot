package types

import "github.com/ivanenkomaksym/remindme_bot/internal/models"

// UserSelectionState holds the complete state of a user's reminder setup
type UserSelectionState struct {
	User            models.User
	RecurrenceType  models.RecurrenceType
	IsWeekly        bool
	WeekOptions     [7]bool
	SelectedTime    string
	ReminderMessage string
	CustomTime      bool
	CustomText      bool
	Language        string
}
