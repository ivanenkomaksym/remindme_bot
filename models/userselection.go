package models

// UserSelection holds the complete state of a user's reminder setup
type UserSelection struct {
	RecurrenceType  RecurrenceType
	IsWeekly        bool
	WeekOptions     [7]bool
	SelectedTime    string
	ReminderMessage string
	CustomTime      bool
	CustomText      bool
}

func CreateUserSelection() UserSelection {
	return UserSelection{
		WeekOptions: [7]bool{false, false, false, false, false, false, false},
	}
}
